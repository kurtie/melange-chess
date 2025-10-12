# Muabdib Chess Engine

Muabdib Chess Engine is an open-source chess engine written in Go. This is my side project. The engine does not aspire to be the fastest (I should have chosen another language for that) or the strongest, but it aims to be a fast Go implementation (I plan to use some x86 assembly).

## Implemented features

- Bitboard representation
- Move generation correcteness validated by perft
- Negamax search tree
- Evaluation function
- UCI protocol

## State

Very early stage. Moves are correctly generated (including castilngs, *en passe* and promotions) and the evaluation function is very basic. UCI protocol implementation is very uncomplete but it is possible to play chess (at least with [nibbler](https://github.com/rooklift/nibbler))