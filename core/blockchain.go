package core

import (
	"fmt"
	"sync"
)

type Blockchain struct {
	store Storage
	//读多写少
	lock sync.RWMutex
	//区块链中块肯定存在，用指针更好
	headers   []*Header
	validator Validator
}

func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(),
	}
	//BlockValidator实现了Validator接口的方法，所以它也属于Validator类型
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

// 创世区块的产生不需要验证
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()
	// logrus.WithFields(logrus.Fields{
	// 	"height": b.Height,
	// 	"hash":   b.Hash(BlockHasher{}),
	// }).Info("adding new block")
	return bc.store.Put(b)
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	//先验证  验证过了就直接添加不用验证
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) getHeader(height uint32) (*Header, error) {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	return bc.headers[height], nil
}
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *Blockchain) Height() uint32 {
	//读锁的延迟释放不会显著影响性能
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	return uint32(len(bc.headers) - 1)
}
