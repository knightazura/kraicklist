# KraickList

Haraj take home challenge!

Simple web app for fictional startup called KraickList. This app will allow users to search ads from given sample data located in `data.gz`.
- First - https://kraicklist-knightazura.herokuapp.com
- Second - https://kraicklist-riot-kaz.herokuapp.com

## Features
- Efficient indexing and search
- Typo tolerance (only for first app)
- Support [BM25 correlation](https://en.wikipedia.org/wiki/Okapi_BM25)
- Support token proximity calculation
- Clean architecture design that easy to add any new text searching implementation

## Installation
1. Create a free Algolia account
2. Set environment variables:
    - `ALGOLIA_API_KEY`
    - `ALGOLIA_APP_ID`
    - `SEARCH_ENGINE_ACTIVE`
    - `MEILISEARCH_HOST`, `MEILISEARCH_PORT` (if using Meilisearch)
3. (Optional) Install Docker if want to try using Meilisearch

## Further improvement
- Implement graceful shutdown
- Build the frontend application for better search experience
- Experiment with integrating local library for text searching and external service(s) indexing at the same time