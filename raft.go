package simpleraft

import (
	"sync"
	"time"
)

// At any given time each server is in one of three states:
type State int

const (
	Follower State = iota
	Candidate
	Leader
)

type LogEntry struct {
	Index   int // position of this log in entire logs
	Term    int // the term in which this log is created by the leader
	Command interface{}
}

// Raft implements raft node
type Raft struct {
	mu sync.Mutex

	// current state of the server
	state State

	// === Persistent state on all servers (Updated on stable storage before responding to RPCs) ===
	// latest term server has seen (initialized to 0 on first boot, increases monotonically)
	currentTerm int
	// -1 means voted for no one
	votedFor int
	log      []LogEntry

	heartbeatTime time.Time
}

type RequestVoteArgs struct {
	Term         int // candidate's term
	CandidateID  int // candidate requesting vote
	LastLogIndex int // index of candidate's last log entry
	LastLogTerm  int // term of candidate's last log entry
}

type RequestVoteReply struct {
	Term        int  // current term, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

// TODO: not complete yet
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	reply.Term = rf.currentTerm
	reply.VoteGranted = false
	if args.Term < rf.currentTerm {
		return
	}

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.state = Follower
		rf.votedFor = -1
	}

	if rf.votedFor == -1 || rf.votedFor == args.CandidateID {
		rf.votedFor = args.CandidateID
		reply.VoteGranted = true
		rf.heartbeatTime = time.Now()
	}

	return
}

type AppendEntriesArgs struct {
	Term         int // leader's term
	LeaderID     int // so follower can redirect clients
	PrevLogIndex int // index of log entry immediately preceding new ones
	PrevLogTerm  int // term of prevLogIndex entry
	// log entries to store (empty for heart beat)
	// may send more than one for efficiency
	Entries      []LogEntry
	LeaderCommit int // leader's commitIndex
}

type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

// TODO: not complete yet
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	reply.Term = rf.currentTerm
	reply.Success = false
	if args.Term < rf.currentTerm {
		return
	}

	if args.Term > rf.currentTerm {
		reply.Term = args.Term
		rf.currentTerm = args.Term
	}
}
