# Notes about streaming downloading


* We want backfetching enabled. When serving check this is set to true
* We could merge allPendingKeyframe/Chunk in one big structure, to tell client that we have all info available. Otherwise we can mimick the spectator server and only provide latest 5  keyframe / 8 chunkId
* We may should then create a mapping of keyFrame <-> chunkId and save it to json alongside all our data. it could be easily done from metadata info reading, and chunkInfo reading
* we need to save the date chunk and keyframe are received
* we need to save the endStartup and startGame



When serving :
* We have to compute, from the data we have the first keyFrame / chunk info available, or endStart
* We may like to think about a start mechanism, so a replay starts with the beginning info

