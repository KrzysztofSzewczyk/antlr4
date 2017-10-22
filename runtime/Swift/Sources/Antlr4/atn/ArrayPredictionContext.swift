/// 
/// Copyright (c) 2012-2017 The ANTLR Project. All rights reserved.
/// Use of this file is governed by the BSD 3-clause license that
/// can be found in the LICENSE.txt file in the project root.
/// 


public class ArrayPredictionContext: PredictionContext {
    /// 
    /// Parent can be null only if full ctx mode and we make an array
    /// from _#EMPTY_ and non-empty. We merge _#EMPTY_ by using null parent and
    /// returnState == _#EMPTY_RETURN_STATE_.
    /// 
    public final var parents: [PredictionContext?]

    /// 
    /// Sorted for merge, no duplicates; if present,
    /// _#EMPTY_RETURN_STATE_ is always last.
    /// 
    public final let returnStates: [Int]

    public convenience init(_ a: SingletonPredictionContext) {
        let parents = [a.parent]
        self.init(parents, [a.returnState])
    }

    public init(_ parents: [PredictionContext?], _ returnStates: [Int]) {

        self.parents = parents
        self.returnStates = returnStates
        super.init(PredictionContext.calculateHashCode(parents, returnStates))
    }

    override
    final public func isEmpty() -> Bool {
        // since EMPTY_RETURN_STATE can only appear in the last position, we
        // don't need to verify that size==1
        return returnStates[0] == PredictionContext.EMPTY_RETURN_STATE
    }

    override
    final public func size() -> Int {
        return returnStates.count
    }

    override
    final public func getParent(_ index: Int) -> PredictionContext? {
        return parents[index]
    }

    override
    final public func getReturnState(_ index: Int) -> Int {
        return returnStates[index]
    }

    override
    public var description: String {
        if isEmpty() {
            return "[]"
        }
        let buf: StringBuilder = StringBuilder()
        buf.append("[")
        let length = returnStates.count

        for i in 0..<length {
            if i > 0 {
                buf.append(", ")
            }
            if returnStates[i] == PredictionContext.EMPTY_RETURN_STATE {
                buf.append("$")
                continue
            }
            buf.append(returnStates[i])
            if parents[i] != nil {
                buf.append(" ")
                buf.append(parents[i].debugDescription)
            } else {
                buf.append("null")
            }
        }
        buf.append("]")
        return buf.toString()
    }

    internal final func combineCommonParents() {

        let length = parents.count
        var uniqueParents: Dictionary<PredictionContext, PredictionContext> =
        Dictionary<PredictionContext, PredictionContext>()
        for p in 0..<length {
            if let parent: PredictionContext = parents[p] {
                // if !uniqueParents.keys.contains(parent) {
                if uniqueParents[parent] == nil {
                    uniqueParents[parent] = parent  // don't replace
                }
            }
        }

        for p in 0..<length {
            if let parent: PredictionContext = parents[p] {
                parents[p] = uniqueParents[parent]
            }
        }

    }
}


public func ==(lhs: ArrayPredictionContext, rhs: ArrayPredictionContext) -> Bool {
    if lhs === rhs {
        return true
    }
    if lhs.hashValue != rhs.hashValue {
        return false
    }

    return lhs.returnStates == rhs.returnStates && lhs.parents == rhs.parents
}

