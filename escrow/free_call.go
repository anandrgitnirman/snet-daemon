package escrow

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type lockingFreeCallUserService struct {
	storage *FreeCallUserStorage
	locker  Locker
	//todo..see if validator is required here
	validator      *FreeCallPaymentValidator
	replicaGroupID func() ([32]byte, error)
}

func NewFreeCallUserService(
	storage *FreeCallUserStorage,

	locker Locker,
	freeCallValidator *FreeCallPaymentValidator, groupIdReader func() ([32]byte, error)) FreeCallUserService {

	return &lockingFreeCallUserService{
		storage:        storage,
		locker:         locker,
		validator:      freeCallValidator,
		replicaGroupID: groupIdReader,
	}
}

func (h *lockingFreeCallUserService) FreeCallUserUsage(key *FreeCallUserKey) (freeCallUser *FreeCallUserData, ok bool, err error) {
	freeCallUser, ok, err = h.storage.Get(key)
	if err != nil {
		return
	}
	if !ok {
		groupId, err := h.replicaGroupID()
		if err != nil {
			return
		}
		return &FreeCallUserData{UserId: key.userId, OrgId: key.organizationId, ServiceId: key.serviceId, GroupID: groupId, FreeCallsMade: 0}, true, nil
	}
	return
}

func (h *lockingFreeCallUserService) ListFreeCallUsers() (users []*FreeCallUserData, err error) {
	return h.storage.GetAll()
}

type freeCallTransaction struct {
	payment      FreeCallPayment
	freeCallUser *FreeCallUserData
	service      *lockingFreeCallUserService
	lock         Lock
}

func (transaction *freeCallTransaction) String() string {
	return fmt.Sprintf("{FreeCallPayment: %v, FreeCallUser: %v}", transaction.payment, transaction.freeCallUser)
}

func (transaction *freeCallTransaction) FreeCallUser() *FreeCallUserData {
	return transaction.freeCallUser
}

func (h *lockingFreeCallUserService) StartFreeCallUserTransaction(payment *FreeCallPayment) (transaction FreeCallTransaction, err error) {
	groupId, err := h.replicaGroupID()
	if err != nil {
		return nil, NewPaymentError(Internal, "cannot get mutex for user: %v", payment.UserId)
	}
	userKey := &FreeCallUserKey{userId: payment.UserId, organizationId: payment.OrganizationId,
		serviceId: payment.ServiceId, groupID: groupId}

	lock, ok, err := h.locker.Lock(userKey.String())
	if err != nil {
		return nil, NewPaymentError(Internal, "cannot get mutex for user: %v", userKey)
	}
	if !ok {
		return nil, NewPaymentError(FailedPrecondition, "another transaction on this user: %v is in progress", userKey)
	}
	defer func(lock Lock) {
		if err != nil {
			e := lock.Unlock()
			if e != nil {
				log.WithError(e).WithField("userKey", userKey).WithField("err", err).Error("Transaction is cancelled because of err, but freeCallUserData cannot be unlocked. All other transactions on this freeCallUserData will be blocked until unlock. Please unlock freeCallUserData manually.")
			}
		}
	}(lock)

	freeCallUserData, ok, err := h.FreeCallUserUsage(userKey)
	if err != nil {
		return nil, NewPaymentError(Internal, "payment freeCallUserData error:"+err.Error())
	}
	if !ok {
		log.Warn("Payment freeCallUserData not found")
		return nil, NewPaymentError(Unauthenticated, "payment freeCallUserData \"%v\" not found", userKey)
	}

	err = h.validator.Validate(payment)
	if err != nil {
		return
	}

	return &freeCallTransaction{
		payment:      *payment,
		freeCallUser: freeCallUserData,
		lock:         lock,
		service:      h,
	}, nil
}

func (transaction *freeCallTransaction) Commit() error {
	defer func(payment *freeCallTransaction) {
		err := payment.lock.Unlock()
		if err != nil {
			log.WithError(err).WithField("transaction", payment).
				Error("free call user cannot be unlocked because of error." +
					" All other transactions on this channel will be blocked until unlock." +
					" Please unlock user for free calls manually.")
		} else {
			log.Debug("free call user unlocked")
		}
	}(transaction)
	group_id, err := transaction.service.replicaGroupID()
	if err != nil {
		log.WithError(err)
		return err
	}
	freeCallUserKey := &FreeCallUserKey{userId: transaction.payment.UserId, organizationId: transaction.payment.OrganizationId,
		serviceId: transaction.payment.ServiceId, groupID: group_id}
	IncrementFreeCallCount(transaction.FreeCallUser())
	e := transaction.service.storage.Put(
		freeCallUserKey,
		transaction.FreeCallUser(),
	)
	if e != nil {
		log.WithError(e).Error("Unable to store new transaction free call user state")
		return NewPaymentError(Internal, "unable to store new transaction free call user state")
	}

	log.Debug("Free Call Payment completed")
	return nil
}

func (payment *freeCallTransaction) Rollback() error {
	defer func(payment *freeCallTransaction) {
		err := payment.lock.Unlock()
		if err != nil {
			log.WithError(err).WithField("payment", payment).Error("free call user cannot be unlocked because of error. All other transactions on this channel will be blocked until unlock. Please unlock channel manually.")
		} else {
			log.Debug("Free call Payment rolled back, free call user unlocked")
		}
	}(payment)
	return nil
}
