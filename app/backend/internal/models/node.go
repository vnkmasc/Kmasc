package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Node của Merkle Tree
type Node struct {
	Hash  string
	Left  *Node
	Right *Node
}

// ProofNode: dùng để lưu MerkleProof cho từng leaf
type ProofNode struct {
	Hash     string `bson:"hash" json:"hash"`
	Position string `bson:"position" json:"position"` // "L" hoặc "R"
}

// MerkleTree struct (có Root)
type MerkleTree struct {
	Root *Node
}

// RootHash trả về root hash của cây
func (t *MerkleTree) RootHash() string {
	if t == nil || t.Root == nil {
		return ""
	}
	return t.Root.Hash
}

// GetProof trả về Merkle proof cho một leaf hash
func (t *MerkleTree) GetProof(target string) []ProofNode {
	if t == nil || t.Root == nil {
		return nil
	}
	return getMerkleProof(t.Root, target)
}

// NewMerkleTree: tạo Merkle Tree từ [][]byte
func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{nil}
	}

	var nodes []*Node
	for _, d := range data {
		hash := sha256.Sum256(d)
		nodes = append(nodes, &Node{Hash: hex.EncodeToString(hash[:])})
	}

	root := buildTree(nodes)
	return &MerkleTree{Root: root}
}

// NewMerkleTreeFromStrings: tạo MerkleTree từ []string (giữ nguyên chuỗi hash)
func NewMerkleTreeFromStrings(data []string) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{nil}
	}

	var nodes []*Node
	for _, s := range data {
		nodes = append(nodes, &Node{Hash: s})
	}

	root := buildTree(nodes)
	return &MerkleTree{Root: root}
}

// buildTree: xây cây từ leaf nodes
func buildTree(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	for len(nodes) > 1 {
		var newLevel []*Node
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				combinedHash := hashConcat(nodes[i].Hash, nodes[i+1].Hash)
				newNode := &Node{
					Hash:  combinedHash,
					Left:  nodes[i],
					Right: nodes[i+1],
				}
				newLevel = append(newLevel, newNode)
			} else {
				newLevel = append(newLevel, nodes[i])
			}
		}
		nodes = newLevel
	}
	return nodes[0]
}

// hashConcat: ghép 2 hash và SHA256
func hashConcat(left, right string) string {
	h := sha256.Sum256([]byte(left + right))
	return hex.EncodeToString(h[:])
}

// getMerkleProof: tạo proof cho leaf
func getMerkleProof(node *Node, target string) []ProofNode {
	var proof []ProofNode
	findPath(node, target, &proof)
	return proof
}

// findPath: tìm đường đi từ root -> leaf và lưu proof
func findPath(node *Node, target string, proof *[]ProofNode) bool {
	if node == nil {
		return false
	}
	if node.Left == nil && node.Right == nil {
		return node.Hash == target
	}

	if findPath(node.Left, target, proof) {
		if node.Right != nil {
			*proof = append(*proof, ProofNode{Hash: node.Right.Hash, Position: "R"})
		}
		return true
	}

	if findPath(node.Right, target, proof) {
		if node.Left != nil {
			*proof = append(*proof, ProofNode{Hash: node.Left.Hash, Position: "L"})
		}
		return true
	}

	return false
}

// VerifyProof: kiểm chứng leaf có thuộc về cây root
func VerifyProof(leaf string, proof []ProofNode, root string) bool {
	hash := leaf
	for _, p := range proof {
		if p.Position == "L" {
			hash = hashConcat(p.Hash, hash)
		} else {
			hash = hashConcat(hash, p.Hash)
		}
	}
	return hash == root
}

// HashEDiplomaInfo: ví dụ hash thông tin diploma (có thể dùng trong service)
func HashEDiplomaInfo(ed *EDiploma) string {
	data := ed.StudentCode + ed.Name + ed.CertificateType + ed.Course + ed.FacultyID.Hex()
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// Optional: timestamp helper
func now() time.Time {
	return time.Now()
}
