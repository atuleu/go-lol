package lol

// QueueID uniquely represents a Queue
type QueueID uint64

const (
	// CUSTOM represents Custom games
	CUSTOM QueueID = 0
	// NORMAL5x5BLIND represents Normal 5v5 Blind Pick games
	NORMAL5x5BLIND QueueID = 2
	// BOT5x5 represents Historical Summoner's Rift Coop vs AI games
	BOT5x5 QueueID = 7
	// BOT5x5INTRO represents Summoner's Rift Coop vs AI Intro Bot
	// games
	BOT5x5INTRO QueueID = 31
	// BOT5x5BEGINNER represents Summoner's Rift Coop vs AI Beginner
	// Bot games
	BOT5x5BEGINNER QueueID = 32
	// BOT5x5INTERMEDIATE represents Historical Summoner's Rift Coop
	// vs AI Intermediate Bot games
	BOT5x5INTERMEDIATE QueueID = 33
	// NORMAL3x3 represents Normal 3v3 games
	NORMAL3x3 QueueID = 8
	// NORMAL5x5DRAFT represents Normal 5v5 Draft Pick games
	NORMAL5x5DRAFT QueueID = 14
	// ODIN5x5BLIND represents Dominion 5v5 Blind Pick games
	ODIN5x5BLIND QueueID = 16
	// ODIN5x5DRAFT represents Dominion 5v5 Draft Pick games
	ODIN5x5DRAFT QueueID = 17
	// BOTODIN5x5 represents Dominion Coop vs AI games
	BOTODIN5x5 QueueID = 25
	// RANKEDSOLO5x5 represents Ranked Solo 5v5 games
	RANKEDSOLO5x5 QueueID = 4
	// RANKEDPREMADE3x3 represents Ranked Premade 3v3 games
	RANKEDPREMADE3x3 QueueID = 9
	// RANKEDPREMADE5x5 represents Ranked Premade 5v5 games
	RANKEDPREMADE5x5 QueueID = 6
	// RANKEDTEAM3x3 represents Ranked Team 3v3 games
	RANKEDTEAM3x3 QueueID = 41
	// RANKEDTEAM5x5 represents Ranked Team 5v5 games
	RANKEDTEAM5x5 QueueID = 42
	// BOTTT3x3 represents Twisted Treeline Coop vs AI games
	BOTTT3x3 QueueID = 52
	// GROUPFINDER5x5 represents Team Builder games
	GROUPFINDER5x5 QueueID = 61
	// ARAM5x5 represents ARAM games
	ARAM5x5 QueueID = 65
	// ONEFORALL5x5 represents One for All games
	ONEFORALL5x5 QueueID = 70
	// FIRSTBLOOD1x1 represents Snowdown Showdown 1v1 games
	FIRSTBLOOD1x1 QueueID = 72
	// FIRSTBLOOD2x2 represents Snowdown Showdown 2v2 games
	FIRSTBLOOD2x2 QueueID = 73
	// SR6x6 represents Summoner's Rift 6x6 Hexakill games
	SR6x6 QueueID = 75
	// URF5x5 represents Ultra Rapid Fire games
	URF5x5 QueueID = 76
	// BOTURF5x5 represents Ultra Rapid Fire games played against AI
	// games
	BOTURF5x5 QueueID = 83
	// NIGHTMAREBOT5x5RANK1 represents Doom Bots Rank 1 games
	NIGHTMAREBOT5x5RANK1 QueueID = 91
	// NIGHTMAREBOT5x5RANK2 represents Doom Bots Rank 2 games
	NIGHTMAREBOT5x5RANK2 QueueID = 92
	// NIGHTMAREBOT5x5RANK5 represents Doom Bots Rank 5 games
	NIGHTMAREBOT5x5RANK5 QueueID = 93
	// ASCENSION5x5 represents Ascension games
	ASCENSION5x5 QueueID = 96
	// HEXAKILL represents Twisted Treeline 6x6 Hexakill games
	HEXAKILL QueueID = 98
	// KINGPORO5x5 represents King Poro games
	KINGPORO5x5 QueueID = 300
	// COUNTERPICK represents Nemesis games
	COUNTERPICK QueueID = 310
)

// A MapID uniquely represents a Map
type MapID uint64

/*const (
	// Original Summer Variant
1 Summoner's Rift
// Original Autumn Variant
2 Summoner's Rift
3 The Proving Grounds Tutorial Map
4 Twisted Treeline Original Version
8 The Crystal Scar Dominion Map
10 Twisted Treeline Current Version
11 Summoner's Rift Current Version
12 Howling Abyss ARAM Map
)*/
