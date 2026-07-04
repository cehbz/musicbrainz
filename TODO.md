# TODO

- [ ] **Build-time `ANALYZE` must run full (`analysis_limit=0`) for stat4.** The pipeline's serving-phase ANALYZE creates `sqlite_stat4` (modernc is STAT4-enabled) but leaves it EMPTY (verified on the live 2026-06-20 mirror, 2026-07-04) — approximate/limited analysis collects no histogram samples. Full stats keep the ingestion pipeline symmetric with the discogs importer's full-ANALYZE step (its stat4 histograms are load-bearing for consumers' query plans). The live mirror got hand-applied full stats 2026-07-04.
- [ ] Monthly rebuild automation (fetch → import → build → verify → swap) — currently manual; tie to the MB fullexport cadence.
