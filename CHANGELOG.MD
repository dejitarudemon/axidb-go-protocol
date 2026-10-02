# Changelog

## 1.0.1 - 2026-10-02

Added:

- Extend on bodies.BatchRequestsBuilder, which appends another batch's requests before the frame is built
- Merge on bodies.BatchResultsBuilder, which lets a BatchAnswer take in other BatchAnswers before the frame is built

## 1.0.0 — 2026-09-26

We fucking did it.

First release of the Go codec for AxiDB.

- Hello (protocol version 0) and working version 1
- Typed values on the wire
- Version chosen per frame after Hello
