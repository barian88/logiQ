package core

import "backend/models"

// NodeKind enumerates supported logical operators.
type NodeKind int

const (
	Var NodeKind = iota
	Not
	And
	Or
	Impl
	Iff
)

// Node represents a propositional formula AST node.
type Node struct {
	Kind        NodeKind
	Name        string
	Left, Right *Node
}

// TruthTablePools holds candidates for truth-table style questions.
type TruthTablePools struct {
	Formula  *Node
	Vars     []string
	TrueSet  []string
	FalseSet []string
}

// EquivalencePools holds equivalence question candidates.
type EquivalencePools struct {
	Target       *Node
	Vars         []string
	EquivPool    []string
	NonEquivPool []string
}

// InferencePools holds inference question candidates.
type InferencePools struct {
	Premises           string
	Vars               []string
	ValidConclusions   []string
	InvalidConclusions []string
}

// CandidatePools aggregates category-specific pools.
type CandidatePools struct {
	TruthTable  *TruthTablePools
	Equivalence *EquivalencePools
	Inference   *InferencePools
}

// Blueprint captures metadata for regenerating a question.
type Blueprint struct {
	Seed           int64
	Difficulty     models.QuestionDifficulty
	Category       models.QuestionCategory
	QType          models.QuestionType
	Intent         string
	MCCorrectCount int
	ProfileSample  map[string]any
	TemplateID     string
	PoolsSummary   map[string]any
}

// Clone returns a deep copy of the node (safe for nil receivers).
func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	clone := &Node{Kind: n.Kind, Name: n.Name}
	clone.Left = n.Left.Clone()
	clone.Right = n.Right.Clone()
	return clone
}
