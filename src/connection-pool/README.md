# A Connection pool for Postgres written in go  
## What is it capable of?
- Initialise a connection pool with given config, up to a max number of connections.
- Maintain separate connections for read and write replica
- Graceful shutdown
- Handle timeouts gracefully 

## Implementation
