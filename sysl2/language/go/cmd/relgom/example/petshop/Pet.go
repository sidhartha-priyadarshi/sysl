//////////////////////////////////////////
//                                      //
//  AUTOGENERATED CODE -- DO NOT EDIT!  //
//                                      //
//////////////////////////////////////////
package petshopmodel

import (
	"encoding/json"
	"time"

	"github.com/anz-bank/sysl/sysl2/language/go/pkg/relgom/relgomlib"
	"github.com/mediocregopher/seq"
)

// petPK is the Key for Pet.
type petPK struct {
	petID int64
}

func (k petPK) Hash(i uint32) uint32 {
	return relgomlib.Hash(i, k.petID)
}

func (k petPK) Equal(i interface{}) bool {
	if l, ok := i.(petPK); ok {
		return (k == l)
	}
	return false
}

// petData is the internal representation of a tuple in the model.
type petData struct {
	petPK
	breedID *int64
	name    *string
	dob     *time.Time
	numLegs *int64
	desexed *bool
}

// MarshalJSON implements json.Marshaler.
func (d *petData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PetID   int64      `json:"petId,omitempty"`
		BreedID *int64     `json:"breedId,omitempty"`
		Name    *string    `json:"name,omitempty"`
		Dob     *time.Time `json:"dob,omitempty"`
		NumLegs *int64     `json:"numLegs,omitempty"`
		Desexed *bool      `json:"desexed,omitempty"`
	}{PetID: d.petID, BreedID: d.breedID, Name: d.name, Dob: d.dob, NumLegs: d.numLegs, Desexed: d.desexed})
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *petData) UnmarshalJSON(data []byte) error {
	var u struct {
		PetID   int64      `json:"petId,omitempty"`
		BreedID *int64     `json:"breedId,omitempty"`
		Name    *string    `json:"name,omitempty"`
		Dob     *time.Time `json:"dob,omitempty"`
		NumLegs *int64     `json:"numLegs,omitempty"`
		Desexed *bool      `json:"desexed,omitempty"`
	}
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	*d = petData{petPK: petPK{petID: u.PetID}, breedID: u.BreedID, name: u.Name, dob: u.Dob, numLegs: u.NumLegs, desexed: u.Desexed}
	return nil
}

// Pet is the public representation tuple in the model.
type Pet struct {
	*petData
	model PetShopModel
}

// PetID gets the petId attribute from the Pet.
func (t Pet) PetID() int64 {
	return t.petID
}

// Breed gets the Breed corresponding to the breedId attribute from t.
func (t Pet) Breed() Breed {
	u, _ := t.model.GetBreed().Lookup(*t.breedID)
	return u
}

// Name gets the name attribute from the Pet.
func (t Pet) Name() *string {
	return t.name
}

// Dob gets the dob attribute from the Pet.
func (t Pet) Dob() *time.Time {
	return t.dob
}

// NumLegs gets the numLegs attribute from the Pet.
func (t Pet) NumLegs() *int64 {
	return t.numLegs
}

// Desexed gets the desexed attribute from the Pet.
func (t Pet) Desexed() *bool {
	return t.desexed
}

// PetBuilder builds an instance of Pet in the model.
type PetBuilder struct {
	petData
	model PetShopModel
	mask  [1]uint64
	apply func(t *petData) (*seq.HashMap, error)
}

// WithBreedID sets the breedId attribute of the PetBuilder from t.
func (b *PetBuilder) WithBreed(t Breed) *PetBuilder {
	relgomlib.UpdateMaskForFieldButPanicIfAlreadySet(&b.mask[0], (uint64(1) << 1))
	b.breedID = &t.breedID
	return b
}

// WithName sets the name attribute of the PetBuilder.
func (b *PetBuilder) WithName(value string) *PetBuilder {
	relgomlib.UpdateMaskForFieldButPanicIfAlreadySet(&b.mask[0], (uint64(1) << 2))
	b.name = &value
	return b
}

// WithDob sets the dob attribute of the PetBuilder.
func (b *PetBuilder) WithDob(value time.Time) *PetBuilder {
	relgomlib.UpdateMaskForFieldButPanicIfAlreadySet(&b.mask[0], (uint64(1) << 3))
	b.dob = &value
	return b
}

// WithNumLegs sets the numLegs attribute of the PetBuilder.
func (b *PetBuilder) WithNumLegs(value int64) *PetBuilder {
	relgomlib.UpdateMaskForFieldButPanicIfAlreadySet(&b.mask[0], (uint64(1) << 4))
	b.numLegs = &value
	return b
}

// WithDesexed sets the desexed attribute of the PetBuilder.
func (b *PetBuilder) WithDesexed(value bool) *PetBuilder {
	relgomlib.UpdateMaskForFieldButPanicIfAlreadySet(&b.mask[0], (uint64(1) << 5))
	b.desexed = &value
	return b
}

var petStaticMetadata = &relgomlib.EntityTypeStaticMetadata{PKMask: []uint64{0x1}, RequiredMask: []uint64{0x0}}

// Apply applies the built Pet.
func (b *PetBuilder) Apply() (PetShopModel, Pet, error) {
	relgomlib.PanicIfRequiredFieldsNotSet(b.mask[:], petStaticMetadata.RequiredMask, ",,,,,")
	set, err := b.apply(&b.petData)
	if err != nil {
		return PetShopModel{}, Pet{}, err
	}
	model, _ := b.model.relations.Set(petKey, petRelationData{set})
	return PetShopModel{model}, Pet{&b.petData, b.model}, nil
}

// petRelationData represents a set of Pet.
type petRelationData struct {
	set *seq.HashMap
}

// Size returns the number of tuples in d.
func (d petRelationData) Size() uint64 {
	return d.set.Size()
}

// MarshalJSON implements json.Marshaler.
func (d petRelationData) MarshalJSON() ([]byte, error) {
	a := make([]*petData, 0, d.set.Size())
	for kv, m, has := d.set.FirstRestKV(); has; kv, m, has = m.FirstRestKV() {
		a = append(a, kv.Val.(*petData))
	}
	return json.Marshal(a)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *petRelationData) UnmarshalJSON(data []byte) error {
	a := []*petData{}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	set := seq.NewHashMap()
	for _, e := range a {
		set, _ = set.Set(e.petPK, e)
	}
	*d = petRelationData{set}
	return nil
}

// PetRelation represents a set of Pet.
type PetRelation struct {
	petRelationData
	model PetShopModel
}

// Insert creates a builder to insert a new Pet.
func (r PetRelation) Insert() *PetBuilder {
	model, id := r.model.newID()
	return &PetBuilder{model: model, apply: func(t *petData) (*seq.HashMap, error) {
		t.petID = int64(id)
		set, _ := r.model.GetPet().set.Set(t.petPK, t)
		return set, nil
	}}
}

// Update creates a builder to update t in the model.
func (r PetRelation) Update(t Pet) *PetBuilder {
	b := &PetBuilder{petData: *t.petData, model: r.model, apply: func(t *petData) (*seq.HashMap, error) {
		set, _ := r.model.GetPet().set.Set(t.petPK, t)
		return set, nil
	}}
	copy(b.mask[:], petStaticMetadata.PKMask)
	return b
}

// Delete deletes t from the model.
func (r PetRelation) Delete(t Pet) (PetShopModel, error) {
	set, _ := r.model.GetPet().set.Del(t.petPK)
	relations, _ := r.model.relations.Set(petKey, petRelationData{set: set})
	return PetShopModel{relations: relations}, nil
}

// Lookup searches Pet by primary key.
func (r PetRelation) Lookup(petID int64) (Pet, bool) {
	if t, has := r.set.Get(petPK{petID: petID}); has {
		return Pet{petData: t.(*petData), model: r.model}, true
	}
	return Pet{}, false
}

// Delete deletes t from the model.
func (r PetRelation) DeleteWhere(where func(t Pet) bool) (PetShopModel, error) {
	model := r.model
	for i := r.Iterator(); i.MoveNext(); {
		t := i.Current()
		if where(t) {
			var err error
			if model, err = model.GetPet().Delete(t); err != nil {
				return PetShopModel{}, err
			}
		}
	}
	return model, nil
}

// Iterator returns an iterator over Pet tuples in r.
func (r PetRelation) Iterator() PetIterator {
	return &petIterator{model: r.model, set: r.set}
}

// petIterator provides for iteration over a set of petIterator tuples.
type PetIterator interface {
	MoveNext() bool
	Current() Pet
}

type petIterator struct {
	model PetShopModel
	set   *seq.HashMap
	t     *Pet
}

// MoveNext implements seq.Setable.
func (i *petIterator) MoveNext() bool {
	kv, set, has := i.set.FirstRestKV()
	if has {
		i.set = set
		i.t = &Pet{petData: kv.Val.(*petData), model: i.model}
	}
	return has
}

// Current implements seq.Setable.
func (i *petIterator) Current() Pet {
	if i.t == nil {
		panic("no current Pet")
	}
	return *i.t
}