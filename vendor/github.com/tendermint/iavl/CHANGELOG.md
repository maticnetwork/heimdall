# Changelog

## 0.8.1

*July 1st, 2018*

BUG FIXES

- fix bug in iterator going outside its range

## 0.8.0 (June 24, 2018)

BREAKING CHANGES

- Nodes are encoded using proto3/amino style integers and byte slices (ie. varints and
  varint prefixed byte slices)
- Unified RangeProof
- Proofs are encoded using Amino
- Hash function changed from RIPEMD160 to the first 20 bytes of SHA256 output

## 0.7.0 (March 21, 2018)

BREAKING CHANGES

- LoadVersion and Load return the loaded version number
    - NOTE: this behaviour was lost previously and we failed to document in changelog,
        but now it's back :)

## 0.6.1 (March 2, 2018)

IMPROVEMENT

- Remove spurious print statement from LoadVersion

## 0.6.0 (March 2, 2018)

BREAKING CHANGES

- NewTree order of arguments swapped
- int -> int64, uint64 -> int64
- NewNode takes a version
- Node serialization format changed so version is written right after size
- SaveVersion takes no args (auto increments)
- tree.Get -> tree.Get64
- nodeDB.SaveBranch does not take a callback
- orphaningTree.SaveVersion -> SaveAs
- proofInnerNode includes Version
- ReadKeyXxxProof consolidated into ReadKeyProof
- KeyAbsentProof doesn't include Version
- KeyRangeProof.Version -> Versions

FEATURES

- Implement chunking algorithm to serialize entire tree

## 0.5.0 (October 27, 2017)

First versioned release!
(Originally accidentally released as v0.2.0)

