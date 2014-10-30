/**
 * Package disjointset implements a DisjointSet (Union-Find) datastructure
 * and an alternative linked list based version.
 */
package disjointset

/**
 * Node for DisjointSet datastructure.
 * Contains the parent id, the node rank and the number
 * of components in the tree (only useful for root nodes)
 */
type DisjointSetNode struct {
	parent int
	rank   int
	size   int
}

/**
 * DisjointSet type, contains the elements (nodes) that
 * the disjoint set has and the number of components that
 * it's representing
 */
type DisjointSet struct {
	elements        []DisjointSetNode
	totalComponents int
}

/**
 * Instantiates a new DisjointSet with 'size' elements
 */
func New(size int) *DisjointSet {
	set := new(DisjointSet)
	set.elements = make([]DisjointSetNode, size, size)
	set.totalComponents = size
	for i := 0; i < size; i++ {
		set.elements[i].parent = i
		set.elements[i].rank = 0
		set.elements[i].size = 1
	}
	return set
}

/**
 * Returns the total number of elements that the DisjointSet set has
 */
func (set *DisjointSet) TotalElements() int {
	return len(set.elements)
}

/**
 * Returns the id of the component to which the node i belongs
 */
func (set *DisjointSet) Find(i int) int {
	if i != set.elements[i].parent {
		set.elements[i].parent = set.Find(set.elements[i].parent)
	}
	return set.elements[i].parent
}

/**
 * Returns the size of the component to which the node p belongs
 */
func (set *DisjointSet) Size(p int) int {
	return set.elements[set.Find(p)].size
}

/**
 * Returns the total number of components that the DisjointSet set has
 */
func (set *DisjointSet) Components() int {
	return set.totalComponents
}

/**
 * Returns true if both nodes p and q belong to the same component
 */
func (set *DisjointSet) Connected(p, q int) bool {
	return set.Find(p) == set.Find(q)
}

/**
 * Merges the two components to which p and q belong. It does nothing
 * if they belong to the same component
 */
func (set *DisjointSet) Union(p, q int) int {
	i := set.Find(p)
	j := set.Find(q)
	if i == j {
		return i
	}

	set.totalComponents--
	if set.elements[i].rank < set.elements[j].rank {
		set.elements[i].parent = j
		set.elements[j].size += set.elements[i].size
		return j
	} else {
		set.elements[j].parent = i
		set.elements[i].size += set.elements[j].size
		if set.elements[i].rank == set.elements[j].rank {
			set.elements[i].rank++
		}
		return i
	}
}
