package antlr4

import (
	"fmt"
)

type IPredictionContext interface {
	hashString() string
	getParent(int) IPredictionContext
	getReturnState(int) int
	equals(IPredictionContext) bool
	length() int
	isEmpty() bool
}

type PredictionContext struct {
	cachedHashString string
}

func NewPredictionContext(cachedHashString string) *PredictionContext {

	pc := new(PredictionContext)
	pc.cachedHashString = cachedHashString

	return pc
}

// Represents {@code $} in local context prediction, which means wildcard.
// {@code//+x =//}.
// /
const (
	PredictionContextEMPTY_RETURN_STATE = 0x7FFFFFFF
)

//var PredictionContextEMPTY *PredictionContext = nil

// Represents {@code $} in an array in full context mode, when {@code $}
// doesn't mean wildcard: {@code $ + x = [$,x]}. Here,
// {@code $} = {@link //EMPTY_RETURN_STATE}.
// /

var PredictionContextglobalNodeCount = 1
var PredictionContextid = PredictionContextglobalNodeCount

// Stores the computed hash code of this {@link PredictionContext}. The hash
// code is computed in parts to match the following reference algorithm.
//
// <pre>
// private int referenceHashCode() {
// int hash = {@link MurmurHash//initialize MurmurHash.initialize}({@link
// //INITIAL_HASH})
//
// for (int i = 0 i &lt {@link //size()} i++) {
// hash = {@link MurmurHash//update MurmurHash.update}(hash, {@link //getParent
// getParent}(i))
// }
//
// for (int i = 0 i &lt {@link //size()} i++) {
// hash = {@link MurmurHash//update MurmurHash.update}(hash, {@link
// //getReturnState getReturnState}(i))
// }
//
// hash = {@link MurmurHash//finish MurmurHash.finish}(hash, 2// {@link
// //size()})
// return hash
// }
// </pre>
// /

// This means only the {@link //EMPTY} context is in set.
func (this *PredictionContext) isEmpty() {
	return this == PredictionContextEMPTY
}

func (this *PredictionContext) hasEmptyPath() {
	return this.getReturnState(this.length() - 1) == PredictionContextEMPTY_RETURN_STATE
}

func (this *PredictionContext) hashString() string {
	return this.cachedHashString
}

func calculateHashString(parent *PredictionContext, returnState int) string {
	return "" + fmt.Sprint(parent) + fmt.Sprint(returnState)
}

func calculateEmptyHashString() string {
	return ""
}

func (this *PredictionContext) getParent(index int) IPredictionContext {
	panic("Not implemented")
}

func (this *PredictionContext) length() int {
	panic("Not implemented")
}

func (this *PredictionContext) getReturnState(index int) int {
	panic("Not implemented")
}

// Used to cache {@link PredictionContext} objects. Its used for the shared
// context cash associated with contexts in DFA states. This cache
// can be used for both lexers and parsers.

type PredictionContextCache struct {
	cache map[IPredictionContext]IPredictionContext
}

func NewPredictionContextCache()  {
	t := new(PredictionContextCache)
	t.cache = make(map[IPredictionContext]IPredictionContext)
	return t
}

// Add a context to the cache and return it. If the context already exists,
// return that one instead and do not add a Newcontext to the cache.
// Protect shared cache from unsafe thread access.
//
func (this *PredictionContextCache) add(ctx IPredictionContext) {
	if (ctx == PredictionContextEMPTY) {
		return PredictionContextEMPTY
	}
	var existing = this.cache[ctx]
	if (existing != nil) {
		return existing
	}
	this.cache[ctx] = ctx
	return ctx
}

func (this *PredictionContextCache) get(ctx IPredictionContext) {
	return this.cache[ctx]
}

func (this *PredictionContextCache) length() int {
	return len(this.cache)
}


type SingletonPredictionContext struct {
	PredictionContext
	parentCtx IPredictionContext
	returnState int
}

func NewSingletonPredictionContext(parent IPredictionContext, returnState int) {
	s := new(SingletonPredictionContext)

	//	var hashString string
	//
	//	if (parent != nil){
	//		hashString = calculateHashString(parent, returnState)
	//	} else {
	//		hashString = calculateEmptyHashString()
	//	}

	panic("Must initializer parent predicition context")
	//	PredictionContext.call(s, hashString)

	s.parentCtx = parent
	s.returnState = returnState

	return s
}

func SingletonPredictionContextcreate(parent IPredictionContext, returnState int) *SingletonPredictionContext {
	if (returnState == PredictionContextEMPTY_RETURN_STATE && parent == nil) {
		// someone can pass in the bits of an array ctx that mean $
		return PredictionContextEMPTY
	} else {
		return NewSingletonPredictionContext(parent, returnState)
	}
}

func (this *SingletonPredictionContext) length() int {
	return 1
}

func (this *SingletonPredictionContext) getParent(index int) IPredictionContext {
	return this.parentCtx
}

func (this *SingletonPredictionContext) getReturnState(index int) int {
	return this.returnState
}

func (this *SingletonPredictionContext) equals(other IPredictionContext) {
	if (this == other) {
		return true
	} else if _, ok := other.(*SingletonPredictionContext); !ok {
		return false
	} else if (this.hashString() != other.hashString()) {
		return false // can't be same if hash is different
	} else {

		otherP := other.(*SingletonPredictionContext)

		if this.returnState != other.getReturnState(0) {
			return false
		} else if(this.parentCtx==nil) {
			return otherP.parentCtx == nil
		} else {
			return this.parentCtx.equals(otherP.parentCtx)
		}
	}
}

func (this *SingletonPredictionContext) hashString() {
	return this.cachedHashString
}

func (this *SingletonPredictionContext) toString() string {
	var up string

	if (this.parentCtx == nil){
		up = ""
	} else {
		up = fmt.Sprint(this.parentCtx)
	}

	if (len(up) == 0) {
		if (this.returnState == PredictionContextEMPTY_RETURN_STATE) {
			return "$"
		} else {
			return "" + this.returnState
		}
	} else {
		return "" + this.returnState + " " + up
	}
}

type EmptyPredictionContext struct {
	SingletonPredictionContext
}

func NewEmptyPredictionContext() *EmptyPredictionContext {

	panic("Must init SingletonPredictionContext")
	//	SingletonPredictionContext.call(this, nil, PredictionContextEMPTY_RETURN_STATE)

	p := new(EmptyPredictionContext)
	return p
}

func (this *EmptyPredictionContext) isEmpty() {
	return true
}

func (this *EmptyPredictionContext) getParent(index int) PredictionContext {
	return nil
}

func (this *EmptyPredictionContext) getReturnState(index int) int {
	return this.returnState
}

func (this *EmptyPredictionContext) equals(other *PredictionContext) bool {
	return this == other
}

func (this *EmptyPredictionContext) toString() string {
	return "$"
}

var PredictionContextEMPTY = NewEmptyPredictionContext()

type ArrayPredictionContext struct {
	PredictionContext
	parents []IPredictionContext
	returnStates []int
}

func NewArrayPredictionContext(parents []IPredictionContext, returnStates []int) *ArrayPredictionContext {
	// Parent can be nil only if full ctx mode and we make an array
	// from {@link //EMPTY} and non-empty. We merge {@link //EMPTY} by using
	// nil parent and
	// returnState == {@link //EMPTY_RETURN_STATE}.

	c := new(ArrayPredictionContext)

	panic("Must init PredictionContext")
	//	var hash = calculateHashString(parents, returnStates)
	//	PredictionContext.call(c, hash)
	c.parents = parents
	c.returnStates = returnStates

	return c
}

func (this *ArrayPredictionContext) isEmpty() {
	// since EMPTY_RETURN_STATE can only appear in the last position, we
	// don't need to verify that size==1
	return this.returnStates[0] == PredictionContextEMPTY_RETURN_STATE
}

func (this *ArrayPredictionContext) length() int {
	return len(this.returnStates)
}

func (this *ArrayPredictionContext) getParent(index int) IPredictionContext {
	return this.parents[index]
}

func (this *ArrayPredictionContext) getReturnState(index int) int {
	return this.returnStates[index]
}

func (this *ArrayPredictionContext) equals(other IPredictionContext) {
	if (this == other) {
		return true
	} else if _, ok := other.(*ArrayPredictionContext); !ok {
		return false
	} else if (this.hashString != other.hashString()) {
		return false // can't be same if hash is different
	} else {
		otherP := other.(*ArrayPredictionContext)
		return this.returnStates == otherP.returnStates && this.parents == otherP.parents
	}
}

func (this *ArrayPredictionContext) toString() string {
	if (this.isEmpty()) {
		return "[]"
	} else {
		var s = "["
		for i := 0; i < len(this.returnStates); i++ {
			if (i > 0) {
				s = s + ", "
			}
			if (this.returnStates[i] == PredictionContextEMPTY_RETURN_STATE) {
				s = s + "$"
				continue
			}
			s = s + this.returnStates[i]
			if (this.parents[i] != nil) {
				s = s + " " + this.parents[i]
			} else {
				s = s + "nil"
			}
		}
		return s + "]"
	}
}

// Convert a {@link RuleContext} tree to a {@link PredictionContext} graph.
// Return {@link //EMPTY} if {@code outerContext} is empty or nil.
// /
func predictionContextFromRuleContext(a *ATN, outerContext *RuleContext) IPredictionContext {
	if (outerContext == nil) {
		outerContext = RuleContextEMPTY
	}
	// if we are in RuleContext of start rule, s, then PredictionContext
	// is EMPTY. Nobody called us. (if we are empty, return empty)
	if (outerContext.parentCtx == nil || outerContext == RuleContextEMPTY) {
		return PredictionContextEMPTY
	}
	// If we have a parent, convert it to a PredictionContext graph
	var parent = predictionContextFromRuleContext(a, outerContext.parentCtx)
	var state = a.states[outerContext.invokingState]
	var transition = state.getTransitions()[0]

	return SingletonPredictionContextcreate(parent, transition.(*RuleTransition).followState.getStateNumber())
}

func calculateListsHashString(parents []PredictionContext, returnStates []int) {
	var s = ""

	for _, p := range parents {
		s += fmt.Sprint(p)
	}

	for _, r := range returnStates {
		s += fmt.Sprint(r)
	}

	return s
}

func merge(a, b IPredictionContext, rootIsWildcard bool, mergeCache *DoubleDict) IPredictionContext {
	// share same graph if both same
	if (a == b) {
		return a
	}

	ac, ok1 := a.(*SingletonPredictionContext)
	bc, ok2 := a.(*SingletonPredictionContext)

	if (ok1 && ok2) {
		return mergeSingletons(ac, bc, rootIsWildcard, mergeCache)
	}
	// At least one of a or b is array
	// If one is $ and rootIsWildcard, return $ as// wildcard
	if (rootIsWildcard) {
		if _, ok := a.(EmptyPredictionContext); ok {
			return a
		}
		if _, ok := b.(EmptyPredictionContext); ok {
			return b
		}
	}
	// convert singleton so both are arrays to normalize
	if _, ok := a.(SingletonPredictionContext); ok {
		a = NewArrayPredictionContext([]IPredictionContext{ a.getParent(0) }, []int{ a.getReturnState(0) })
	}
	if _, ok := b.(SingletonPredictionContext); ok {
		b = NewArrayPredictionContext( []IPredictionContext{ b.getParent(0) }, []int{ b.getReturnState(0) })
	}
	return mergeArrays(a, b, rootIsWildcard, mergeCache)
}

//
// Merge two {@link SingletonPredictionContext} instances.
//
// <p>Stack tops equal, parents merge is same return left graph.<br>
// <embed src="images/SingletonMerge_SameRootSamePar.svg"
// type="image/svg+xml"/></p>
//
// <p>Same stack top, parents differ merge parents giving array node, then
// remainders of those graphs. A Newroot node is created to point to the
// merged parents.<br>
// <embed src="images/SingletonMerge_SameRootDiffPar.svg"
// type="image/svg+xml"/></p>
//
// <p>Different stack tops pointing to same parent. Make array node for the
// root where both element in the root point to the same (original)
// parent.<br>
// <embed src="images/SingletonMerge_DiffRootSamePar.svg"
// type="image/svg+xml"/></p>
//
// <p>Different stack tops pointing to different parents. Make array node for
// the root where each element points to the corresponding original
// parent.<br>
// <embed src="images/SingletonMerge_DiffRootDiffPar.svg"
// type="image/svg+xml"/></p>
//
// @param a the first {@link SingletonPredictionContext}
// @param b the second {@link SingletonPredictionContext}
// @param rootIsWildcard {@code true} if this is a local-context merge,
// otherwise false to indicate a full-context merge
// @param mergeCache
// /
func mergeSingletons(a, b *SingletonPredictionContext, rootIsWildcard bool, mergeCache *DoubleDict) IPredictionContext {
	if (mergeCache != nil) {
		var previous = mergeCache.get(a, b)
		if (previous != nil) {
			return previous
		}
		previous = mergeCache.get(b, a)
		if (previous != nil) {
			return previous
		}
	}

	var rootMerge = mergeRoot(a, b, rootIsWildcard)
	if (rootMerge != nil) {
		if (mergeCache != nil) {
			mergeCache.set(a, b, rootMerge)
		}
		return rootMerge
	}
	if (a.returnState == b.returnState) {
		var parent = merge(a.parentCtx, b.parentCtx, rootIsWildcard, mergeCache)
		// if parent is same as existing a or b parent or reduced to a parent,
		// return it
		if (parent == a.parentCtx) {
			return a // ax + bx = ax, if a=b
		}
		if (parent == b.parentCtx) {
			return b // ax + bx = bx, if a=b
		}
		// else: ax + ay = a'[x,y]
		// merge parents x and y, giving array node with x,y then remainders
		// of those graphs. dup a, a' points at merged array
		// Newjoined parent so create Newsingleton pointing to it, a'
		var spc = SingletonPredictionContextcreate(parent, a.returnState)
		if (mergeCache != nil) {
			mergeCache.set(a, b, spc)
		}
		return spc
	} else { // a != b payloads differ
		// see if we can collapse parents due to $+x parents if local ctx
		var singleParent IPredictionContext = nil
		if (a == b || (a.parentCtx != nil && a.parentCtx == b.parentCtx)) { // ax +
			// bx =
			// [a,b]x
			singleParent = a.parentCtx
		}
		if (singleParent != nil) { // parents are same
			// sort payloads and use same parent
			var payloads = []int{ a.returnState, b.returnState }
			if (a.returnState > b.returnState) {
				payloads[0] = b.returnState
				payloads[1] = a.returnState
			}
			var parents = []IPredictionContext{ singleParent, singleParent }
			var apc = NewArrayPredictionContext(parents, payloads)
			if (mergeCache != nil) {
				mergeCache.set(a, b, apc)
			}
			return apc
		}
		// parents differ and can't merge them. Just pack together
		// into array can't merge.
		// ax + by = [ax,by]
		var payloads = []int{ a.returnState, b.returnState }
		var parents = []IPredictionContext{ a.parentCtx, b.parentCtx }
		if (a.returnState > b.returnState) { // sort by payload
			payloads[0] = b.returnState
			payloads[1] = a.returnState
			parents = []IPredictionContext{  b.parentCtx, a.parentCtx }
		}
		var a_ = NewArrayPredictionContext(parents, payloads)
		if (mergeCache != nil) {
			mergeCache.set(a, b, a_)
		}
		return a_
	}
}

//
// Handle case where at least one of {@code a} or {@code b} is
// {@link //EMPTY}. In the following diagrams, the symbol {@code $} is used
// to represent {@link //EMPTY}.
//
// <h2>Local-Context Merges</h2>
//
// <p>These local-context merge operations are used when {@code rootIsWildcard}
// is true.</p>
//
// <p>{@link //EMPTY} is superset of any graph return {@link //EMPTY}.<br>
// <embed src="images/LocalMerge_EmptyRoot.svg" type="image/svg+xml"/></p>
//
// <p>{@link //EMPTY} and anything is {@code //EMPTY}, so merged parent is
// {@code //EMPTY} return left graph.<br>
// <embed src="images/LocalMerge_EmptyParent.svg" type="image/svg+xml"/></p>
//
// <p>Special case of last merge if local context.<br>
// <embed src="images/LocalMerge_DiffRoots.svg" type="image/svg+xml"/></p>
//
// <h2>Full-Context Merges</h2>
//
// <p>These full-context merge operations are used when {@code rootIsWildcard}
// is false.</p>
//
// <p><embed src="images/FullMerge_EmptyRoots.svg" type="image/svg+xml"/></p>
//
// <p>Must keep all contexts {@link //EMPTY} in array is a special value (and
// nil parent).<br>
// <embed src="images/FullMerge_EmptyRoot.svg" type="image/svg+xml"/></p>
//
// <p><embed src="images/FullMerge_SameRoot.svg" type="image/svg+xml"/></p>
//
// @param a the first {@link SingletonPredictionContext}
// @param b the second {@link SingletonPredictionContext}
// @param rootIsWildcard {@code true} if this is a local-context merge,
// otherwise false to indicate a full-context merge
// /
func mergeRoot(a, b *SingletonPredictionContext, rootIsWildcard bool) IPredictionContext {
	if (rootIsWildcard) {
		if (a == PredictionContextEMPTY) {
			return PredictionContextEMPTY // // + b =//
		}
		if (b == PredictionContextEMPTY) {
			return PredictionContextEMPTY // a +// =//
		}
	} else {
		if (a == PredictionContextEMPTY && b == PredictionContextEMPTY) {
			return PredictionContextEMPTY // $ + $ = $
		} else if (a == PredictionContextEMPTY) { // $ + x = [$,x]
			var payloads = []int{ b.returnState, PredictionContextEMPTY_RETURN_STATE }
			var parents = []IPredictionContext{ b.parentCtx, nil }
			return NewArrayPredictionContext(parents, payloads)
		} else if (b == PredictionContextEMPTY) { // x + $ = [$,x] ($ is always first if present)
			var payloads = []int{ a.returnState, PredictionContextEMPTY_RETURN_STATE }
			var parents = []IPredictionContext{ a.parentCtx, nil }
			return NewArrayPredictionContext(parents, payloads)
		}
	}
	return nil
}

//
// Merge two {@link ArrayPredictionContext} instances.
//
// <p>Different tops, different parents.<br>
// <embed src="images/ArrayMerge_DiffTopDiffPar.svg" type="image/svg+xml"/></p>
//
// <p>Shared top, same parents.<br>
// <embed src="images/ArrayMerge_ShareTopSamePar.svg" type="image/svg+xml"/></p>
//
// <p>Shared top, different parents.<br>
// <embed src="images/ArrayMerge_ShareTopDiffPar.svg" type="image/svg+xml"/></p>
//
// <p>Shared top, all shared parents.<br>
// <embed src="images/ArrayMerge_ShareTopSharePar.svg"
// type="image/svg+xml"/></p>
//
// <p>Equal tops, merge parents and reduce top to
// {@link SingletonPredictionContext}.<br>
// <embed src="images/ArrayMerge_EqualTop.svg" type="image/svg+xml"/></p>
// /
func mergeArrays(a, b *ArrayPredictionContext, rootIsWildcard bool, mergeCache *DoubleDict) IPredictionContext {
	if (mergeCache != nil) {
		var previous = mergeCache.get(a, b)
		if (previous != nil) {
			return previous
		}
		previous = mergeCache.get(b, a)
		if (previous != nil) {
			return previous
		}
	}
	// merge sorted payloads a + b => M
	var i = 0 // walks a
	var j = 0 // walks b
	var k = 0 // walks target M array

	var mergedReturnStates = make([]int,0)
	var mergedParents = make([]IPredictionContext,0)
	// walk and merge to yield mergedParents, mergedReturnStates
	for i < len(a.returnStates) && j < len(b.returnStates) {
		var a_parent = a.parents[i]
		var b_parent = b.parents[j]
		if (a.returnStates[i] == b.returnStates[j]) {
			// same payload (stack tops are equal), must yield merged singleton
			var payload = a.returnStates[i]
			// $+$ = $
			var bothDollars = payload == PredictionContextEMPTY_RETURN_STATE && a_parent == nil && b_parent == nil
			var ax_ax = (a_parent != nil && b_parent != nil && a_parent == b_parent) // ax+ax
			// ->
			// ax
			if (bothDollars || ax_ax) {
				mergedParents[k] = a_parent // choose left
				mergedReturnStates[k] = payload
			} else { // ax+ay -> a'[x,y]
				var mergedParent = merge(a_parent, b_parent, rootIsWildcard, mergeCache)
				mergedParents[k] = mergedParent
				mergedReturnStates[k] = payload
			}
			i += 1 // hop over left one as usual
			j += 1 // but also skip one in right side since we merge
		} else if (a.returnStates[i] < b.returnStates[j]) { // copy a[i] to M
			mergedParents[k] = a_parent
			mergedReturnStates[k] = a.returnStates[i]
			i += 1
		} else { // b > a, copy b[j] to M
			mergedParents[k] = b_parent
			mergedReturnStates[k] = b.returnStates[j]
			j += 1
		}
		k += 1
	}
	// copy over any payloads remaining in either array
	if (i < len(a.returnStates)) {
		for p := i; p < len(a.returnStates); p++ {
			mergedParents[k] = a.parents[p]
			mergedReturnStates[k] = a.returnStates[p]
			k += 1
		}
	} else {
		for p := j; p < len(b.returnStates); p++ {
			mergedParents[k] = b.parents[p]
			mergedReturnStates[k] = b.returnStates[p]
			k += 1
		}
	}
	// trim merged if we combined a few that had same stack tops
	if (k < len(mergedParents)) { // write index < last position trim
		if (k == 1) { // for just one merged element, return singleton top
			var a_ = SingletonPredictionContextcreate(mergedParents[0], mergedReturnStates[0])
			if (mergeCache != nil) {
				mergeCache.set(a, b, a_)
			}
			return a_
		}
		mergedParents = mergedParents[0:k]
		mergedReturnStates = mergedReturnStates[0:k]
	}

	var M = NewArrayPredictionContext(mergedParents, mergedReturnStates)

	// if we created same array as a or b, return that instead
	// TODO: track whether this is possible above during merge sort for speed
	if (M == a) {
		if (mergeCache != nil) {
			mergeCache.set(a, b, a)
		}
		return a
	}
	if (M == b) {
		if (mergeCache != nil) {
			mergeCache.set(a, b, b)
		}
		return b
	}
	combineCommonParents(mergedParents)

	if (mergeCache != nil) {
		mergeCache.set(a, b, M)
	}
	return M
}

//
// Make pass over all <em>M</em> {@code parents} merge any {@code equals()}
// ones.
// /
func combineCommonParents(parents []IPredictionContext) {
	var uniqueParents = make(map[IPredictionContext]IPredictionContext)

	for p := 0; p < len(parents); p++ {
		var parent = parents[p]
		if uniqueParents[parent] == nil {
			uniqueParents[parent] = parent
		}
	}
	for q := 0; q < len(parents); q++ {
		parents[q] = uniqueParents[parents[q]]
	}
}

func getCachedPredictionContext(context IPredictionContext, contextCache *PredictionContextCache, visited map[IPredictionContext]IPredictionContext) IPredictionContext {

	panic("getCachedPredictionContext not implemented")

	return nil
	//	if (context.isEmpty()) {
	//		return context
	//	}
	//	var existing = visited[context] || nil
	//	if (existing != nil) {
	//		return existing
	//	}
	//	existing = contextCache.get(context)
	//	if (existing != nil) {
	//		visited[context] = existing
	//		return existing
	//	}
	//	var changed = false
	//	var parents = []
	//	for i := 0; i < len(parents); i++ {
	//		var parent = getCachedPredictionContext(context.getParent(i), contextCache, visited)
	//		if (changed || parent != context.getParent(i)) {
	//			if (!changed) {
	//				parents = []
	//				for j := 0; j < len(context); j++ {
	//					parents[j] = context.getParent(j)
	//				}
	//				changed = true
	//			}
	//			parents[i] = parent
	//		}
	//	}
	//	if (!changed) {
	//		contextCache.add(context)
	//		visited[context] = context
	//		return context
	//	}
	//	var updated = nil
	//	if (parents.length == 0) {
	//		updated = PredictionContextEMPTY
	//	} else if (parents.length == 1) {
	//		updated = SingletonPredictionContext.create(parents[0], context.getReturnState(0))
	//	} else {
	//		updated = NewArrayPredictionContext(parents, context.returnStates)
	//	}
	//	contextCache.add(updated)
	//	visited[updated] = updated
	//	visited[context] = updated
	//
	//	return updated
}

// ter's recursive version of Sam's getAllNodes()
//func getAllContextNodes(context, nodes, visited) {
//	if (nodes == nil) {
//		nodes = []
//		return getAllContextNodes(context, nodes, visited)
//	} else if (visited == nil) {
//		visited = {}
//		return getAllContextNodes(context, nodes, visited)
//	} else {
//		if (context == nil || visited[context] != nil) {
//			return nodes
//		}
//		visited[context] = context
//		nodes.push(context)
//		for i := 0; i < len(context); i++ {
//			getAllContextNodes(context.getParent(i), nodes, visited)
//		}
//		return nodes
//	}
//}







