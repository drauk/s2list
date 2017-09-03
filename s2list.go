// src/go/s2list.go   2017-9-2   Alan U. Kennington.
// $Id: s2list.go 46551 2017-09-01 04:37:04Z akenning $
// Singly linked first/last-list for first test program for learning "go".
// Using version go1.1.2.
/*-------------------------------------------------------------------------
Functions in this package.

List_node::
List_node::GetNext
List_node::unlink
List_node::SetValue
List_node::GetValue
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
List_base::
List_base::Empty
List_base::GetFirst
List_base::Length
List_base::ValidLength
List_base::Append
List_base::AppendValue
List_base::Prepend
List_base::PrependValue
List_base::Popfirst
List_base::Poplast
List_base::Found
List_base::Remove
List_base::Clear
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
List_iter::
List_iter::Init
List_iter::Restart
List_iter::ItemCount
List_iter::ItemCountValid
List_iter::Next
-------------------------------------------------------------------------*/

/*
The s2list package implements singly-linked lists. Every node has a
parent-pointer to indicate which list it is contained in. This helps to ensure
structural integrity. For example, a node is preventing from being an element of
the same list twice, or an element of two different lists simultaneously.
*/
package s2list

// External libraries.
// import "fmt"
// import "io"
// import "log"
// import "time"
// import "errors"
// import "net/http"

import "github.com/drauk/elist"

//=============================================================================
//=============================================================================

/*
A List_node is an element of a List_base. Every List_node has a base-pointer
which is nil if the node is not in a list, or points to the List_base if it is
within that list. When a List_node is removed from a list, the base-pointer is
set to nil.
    next *List_node // Next node in a singly linked list.
    base *List_base // The base in which this object is listed.
    value interface{} // The payload of the list node.
*/
type List_node struct {
    //----------------------//
    //      List_node::     //
    //----------------------//
    /*------------------------------------------------------------------------------
      The "next" field is nil if the node is not a member of a list.
      The "base" field is nil if the node is not a member of a list,
      and it points to the list if it is a member of it.
      There is almost no control over what goes into the "value" field.
      Any attempt to make a list homogeneous can be easily defeated by a user
      calling List_node::SetValue().
      However, if a client class keeps the list and nodes private, other users
      should not be able to invalidate the homogeneity of the list.
      ------------------------------------------------------------------------------*/
    next *List_node // Next node in a singly linked list.
    base *List_base // The base in which this object is listed.

    value interface{} // The payload of the list node.
}

/*
List_node::GetNext() returns the next-field of the node. This gives read-only
access to address of the next node in the list, if the node is within a list.
The return value is nil if the node is outside any list, or it is the last
element of its container-list.
*/
func (p *List_node) GetNext() (*List_node, error) {
    //----------------------//
    //  List_node::GetNext  //
    //----------------------//
    if p == nil {
        return nil, elist.New("List_node::GetNext: p == nil")
    }
    return p.next, nil
}   // End of function List_node::GetNext.

func (p *List_node) unlink() error {
    //----------------------//
    //   List_node::unlink  //
    //----------------------//
    /*------------------------------------------------------------------------------
      List_node::unlink()
      Private member function for internal use in this package.
      This should only be called when a node is popped/removed/cleared from a list.
      The value payload is unaffected.
      ------------------------------------------------------------------------------*/
    if p == nil {
        return elist.New("List_node::unlink: p == nil")
    }
    p.next = nil
    p.base = nil
    return nil
}   // End of function List_node::unlink.

/*
List_node::SetValue() clobbers whatever was in the "value" field before.
*/
func (p *List_node) SetValue(v interface{}) error {
    //----------------------//
    //  List_node::SetValue //
    //----------------------//
    /*------------------------------------------------------------------------------
      Apparently Go-language does a copy-on-assign here, or something like that.
      So you end up with two distinct unlinked objects, your original value
      and the value which you copy into the "value" field.
      ------------------------------------------------------------------------------*/
    if p == nil {
        return elist.New("List_node::SetValue: p == nil")
    }
    p.value = v
    return nil
}   // End of function List_node::SetValue.

/*
List_node::GetValue() returns the value-field of the node.
*/
func (p *List_node) GetValue() (interface{}, error) {
    //----------------------//
    //  List_node::GetValue //
    //----------------------//
    /*------------------------------------------------------------------------------
      List_node::GetValue()
      Does this return a simple copy of bits in the value member?
      ------------------------------------------------------------------------------*/
    if p == nil {
        return nil, elist.New("List_node::GetValue: p == nil")
    }
    return p.value, nil
}   // End of function List_node::GetValue.

//=============================================================================
//=============================================================================

/*
List_base is a singly-linked list implementation which has pointers to both the
head and tail elements of the list.
    first *List_node // First node of the list.
    last  *List_node // Last node of the list.
Every node in the list has a base-pointer which points to the list-base which it
is contained in, or which equals nil if the node is not contained in a list.
Various checks are made by List_base methods to prevent corruption of the list
structure, but to further protect the integrity of these lists, it is best to
not re-export list bases or list nodes to other packages.
*/
type List_base struct {
    //----------------------//
    //      List_base::     //
    //----------------------//
    /*------------------------------------------------------------------------------
      All elements of the list must point to the base.
      ------------------------------------------------------------------------------*/
    first *List_node // First node of the list.
    last  *List_node // Last node of the list.
}

/*
List_base::Empty() returns true when the list is empty!
(Obviousness is a feature, not a bug.)
*/
func (p *List_base) Empty() bool {
    //----------------------//
    //   List_base::Empty   //
    //----------------------//
    if p == nil {
        return true
    }
    if p.first == nil {
        return true
    }
    return false
}   // End of function List_base::Empty.

/*
List_base::GetFirst() gives read-only access to the address of the first node
in the list. This can be used in conjunction with List_node::GetNext() to
traverse a list.
*/
func (p *List_base) GetFirst() *List_node {
    //----------------------//
    //  List_base::GetFirst //
    //----------------------//
    if p == nil {
        return nil
    }
    return p.first
}   // End of function List_base::GetFirst.

/*
List_base::Length() will return the wrong answer if the number of elements in
the list does not fit into an int. This function does not verify that all
elements of the list have valid base-pointers.
*/
func (p *List_base) Length() int {
    //----------------------//
    //   List_base::Length  //
    //----------------------//
    if p == nil {
        return 0
    }
    var n int = 0
    if p.first != nil && p.last != nil {
        for q := p.first; q != nil; q = q.next {
            n += 1
        }
    }
    return n
}   // End of function List_base::Length.

/*
List_base::ValidLength() will return the wrong answer if the number of elements
in the list does not fit into an int.
This function checks how many elements of the list have valid base-pointers.
    Return values: (nil, wrong, total)
    nil     number of nil base pointers
    wrong   number of wrong non-nil base pointers
    total   number of elements in the list, right or wrong.
The purpose of this function is to test integrity in case there is a suspicion
of structure-corruption.
*/
func (p *List_base) ValidLength() (int, int, int) {
    //--------------------------//
    //  List_base::ValidLength  //
    //--------------------------//
    if p == nil {
        return 0, 0, 0
    }
    var n_nil, n_wrong, n_total int
    if p.first == nil || p.last == nil {
        return 0, 0, 0
    }
    for q := p.first; q != nil; q = q.next {
        n_total += 1
        if q.base == nil {
            n_nil += 1
        } else if q.base != p {
            n_wrong += 1
        }
    }
    return n_nil, n_wrong, n_total
}   // End of function List_base::ValidLength.

/*
List_base::Append() appends a list-node and links it to the base.
The argument must be the address of a heap-allocated List_node.
Never try to simultaneously put a List_node into the same list twice or into two
different lists.
*/
func (p *List_base) Append(pnode *List_node) error {
    //----------------------//
    //   List_base::Append  //
    //----------------------//
    if p == nil {
        return elist.New("List_base::Append: p == nil")
    }
    if pnode == nil {
        return nil
    }
    // Can't put an object in multiple lists.
    if pnode.base != nil {
        return elist.New("List_base::Append: pnode.base != nil")
    }
    pnode.base = p // Register the node with this list-base.
    pnode.next = nil
    if p.last != nil {
        p.last.next = pnode
    } else {
        p.first = pnode
    }
    p.last = pnode
    return nil
}   // End of function List_base::Append.

/*
List_base::AppendValue() copies the given value to a newly created node, appends
the node to the base, and links the node to the base to indicate which list-base
the list-node belongs to.
*/
func (p *List_base) AppendValue(v interface{}) error {
    //--------------------------//
    //  List_base::AppendValue  //
    //--------------------------//
    if p == nil {
        return elist.New("List_base::AppendValue: p == nil")
    }
    var pnode *List_node = new(List_node)
    var E error

    E = pnode.SetValue(v)
    if E != nil {
        return elist.Push(E, "List_base::AppendValue: pnode.SetValue(v)")
    }
    E = p.Append(pnode)
    if E != nil {
        return elist.Push(E, "List_base::AppendValue: p.Append(pnode)")
    }
    return nil
}   // End of function List_base::AppendValue.

/*
List_base::Prepend() prepends a list-node and links it to the base.
The argument must be the address of a heap-allocated List_node.
Never try to simultaneously put a List_node into the same list twice or into two
different lists.
*/
func (p *List_base) Prepend(pnode *List_node) error {
    //----------------------//
    //  List_base::Prepend  //
    //----------------------//
    if p == nil {
        return elist.New("List_base::Prepend: p == nil")
    }
    if pnode == nil {
        return nil
    }
    // Can't put an object in multiple lists.
    if pnode.base != nil {
        return elist.New("List_base::Prepend: pnode.base != nil")
    }
    pnode.base = p // Register the node with this list-base.
    pnode.next = p.first
    p.first = pnode
    if p.last == nil {
        p.last = pnode
    }
    return nil
}   // End of function List_base::Prepend.

/*
List_base::PrependValue() copies the given value to a newly created node,
prepends the node to the base, and links the node to the base to indicate which
list-base the list-node belongs to.
*/
func (p *List_base) PrependValue(v interface{}) error {
    //--------------------------//
    //  List_base::PrependValue //
    //--------------------------//
    if p == nil {
        return elist.New("List_base::PrependValue: p == nil")
    }
    var pnode *List_node = new(List_node)
    var E error

    E = pnode.SetValue(v)
    if E != nil {
        return elist.Push(E, "List_base::PrependValue: pnode.SetValue(v)")
    }
    E = p.Prepend(pnode)
    if E != nil {
        return elist.Push(E, "List_base::PrependValue: p.Prepend(pnode)")
    }
    return nil
}   // End of function List_base::PrependValue.

/*
List_base::Popfirst() pops the first node from the list and returns it to the
caller. If the list is empty, the nil node-pointer is returned and the error
returned is then nil.
*/
func (p *List_base) Popfirst() (*List_node, error) {
    //----------------------//
    //  List_base::Popfirst //
    //----------------------//
    if p == nil {
        return nil, elist.New("List_base::Popfirst: p == nil")
    }
    if p.first == nil {
        return nil, nil
    }
    // If "first" is nil and "last" is not, this is a very serious error!
    if p.last == nil {
        return nil, elist.New("List_base::Popfirst: p.first != p.last == nil")
    }
    if p.last == p.first {
        p.last = nil
    }
    pnode := p.first
    p.first = pnode.next
    pnode.unlink()
    return pnode, nil
}   // End of function List_base::Popfirst.

/*
List_base::Poplast() pops the last node from the list and returns it to the
caller. If the list is empty, the nil node-pointer is returned and the error
returned is then nil.
*/
func (p *List_base) Poplast() (*List_node, error) {
    //----------------------//
    //  List_base::Poplast  //
    //----------------------//
    if p == nil {
        return nil, elist.New("List_base::Poplast: p == nil")
    }
    if p.first == nil {
        return nil, nil
    }
    // List integrity check.
    // If "first" is nil and "last" is not, the list is corrupted.
    if p.last == nil {
        return nil, elist.New("List_base::Poplast: p.first != p.last == nil")
    }
    var pnode *List_node = nil
    // Special case of only one item found in the list.
    if p.last == p.first {
        pnode = p.first
        p.first = nil
        p.last = nil
        pnode.unlink()
        return pnode, nil
    }
    // Find the second-to-last item in the list.
    var q *List_node
    for q = p.first; q != nil; q = q.next {
        if q.next == p.last {
            break
        }
    }
    // This should never happen. Indicates list is corrupted.
    if q == nil {
        return nil, elist.New("List_base::Poplast: q == nil")
    }
    pnode = p.last
    q.next = nil
    p.last = q
    pnode.unlink()
    return pnode, nil
}   // End of function List_base::Poplast.

/*
List_base::Found(*List_node) returns true if and only if the node is currently
contained in the list.
Neither the node nor the list are affected. This is a read-only method.
*/
func (p *List_base) Found(q *List_node) (bool, error) {
    //----------------------//
    //   List_base::Found   //
    //----------------------//
    if p == nil {
        return false, elist.New("List_base::Found: p == nil")
    }
    // Can't find a nil object in any list.
    if q == nil {
        return false, nil
    }
    // Can't find a non-nil object in an empty list.
    if p.first == nil {
        return false, nil
    }
    // List integrity check.
    // If "first" is nil and "last" is not, this is a very serious error!
    if p.last == nil {
        return false, elist.New("List_base::Found: p.first != p.last == nil")
    }
    // The given object does not belong to this list. So don't even try.
    if q.base != p {
        return false, elist.New("List_base::Found: q.base != p")
    }
    // Try to find q in the list.
    for pnode := p.first; pnode != nil; pnode = pnode.next {
        if pnode == q {
            return true, nil
        }
    }
    return false, nil
}   // End of function List_base::Found.

/*
List_base::Remove() removes the given node from the list, if it is a valid
member of the list, and returns the removed node to the caller.
The returned node is always either nil or the same as the requested node.
*/
func (p *List_base) Remove(q *List_node) (*List_node, error) {
    //----------------------//
    //   List_base::Remove  //
    //----------------------//
    if p == nil {
        return nil, elist.New("List_base::Remove: p == nil")
    }
    // Can't find a nil object in any list.
    if q == nil {
        return nil, nil
    }
    // Can't find a non-nil object in an empty list.
    if p.first == nil {
        return nil, nil
    }
    // List integrity check.
    // If "first" is nil and "last" is not, this is a very serious error!
    if p.last == nil {
        return nil, elist.New("List_base::Remove: p.first != p.last == nil")
    }
    // The given object does not belong to the list.
    if q.base != p {
        return nil, elist.New("List_base::Remove: q.base != p")
    }
    // Special case of popping the first element.
    if p.first == q {
        if p.last == p.first {
            p.last = nil
        }
        p.first = q.next

        // Unlink the node from the list base.
        q.unlink()
        return q, nil
    }
    // Try to find the predecessor of q in the list.
    var pnode *List_node
    for pnode = p.first; pnode != nil; pnode = pnode.next {
        if pnode.next == q {
            break
        }
    }
    // Didn't find the object in the list. Should never happen!
    if pnode == nil {
        return nil, elist.New("List_base::Remove: pnode == nil")
    }
    pnode.next = q.next
    if p.last == q {
        p.last = pnode
    }
    // Unlink the node from the list.
    q.unlink()
    return q, nil
}   // End of function List_base::Remove.

/*
List_base::Clear() removes all nodes from the list and casts them adrift.
There is no destructor in the Go-language.
However, the base-pointers of the removed nodes are set to nil to indicate their
cast-adrift status.
(This function does not check that the nodes in the last have the correct
base-pointer.)
*/
func (p *List_base) Clear() error {
    //----------------------//
    //   List_base::Clear   //
    //----------------------//
    if p == nil {
        return elist.New("List_base::Clear: p == nil")
    }
    if p.first == nil {
        return nil
    }
    // If "first" is nil and "last" is not, this is a very serious error!
    if p.last == nil {
        return elist.New("List_base::Clear: p.first != p.last == nil")
    }
    // Pop and unlink the first element recursively until nothing is left.
    for p.first != nil {
        if p.last == p.first {
            p.last = nil
        }
        pnode := p.first
        p.first = pnode.next
        pnode.unlink()
    }
    return nil
}   // End of function List_base::Clear.

//=============================================================================
//=============================================================================

/*
List_iter is used for traversals of the nodes in a List_base.
    base    *List_base // The list which is used for the iteration.
    current *List_node // The last node delivered by the iterator.
List_base lists can also be traversed using the List_base::GetFirst() and
List_node::GetNext() functions. However, the List_iter::Next() function performs
integrity checks to ensure valid results. (For example, a node could be moved
from one list to another between List_node::GetNext() calls, which would result
in the traversal continuing from the original list to a different list!)
*/
type List_iter struct {
    //----------------------//
    //      List_iter::     //
    //----------------------//
    base    *List_base // The list which is used for the iteration.
    current *List_node // The last node delivered by the iterator.
}

/*
List_iter::Init() initializes a list-iterator to point at a given list-base.
*/
func (p *List_iter) Init(b *List_base) error {
    //----------------------//
    //    List_iter::Init   //
    //----------------------//
    if p == nil {
        return elist.New("List_base::Init: p == nil")
    }
    p.base = b
    p.current = nil
    return nil
}   // End of function List_iter::Init.

/*
List_iter::Restart() rewinds the current node-pointer to the start of the list.
*/
func (p *List_iter) Restart() error {
    //----------------------//
    //  List_iter::Restart  //
    //----------------------//
    if p == nil {
        return elist.New("List_base::Restart: p == nil")
    }
    p.current = nil
    return nil
}   // End of function List_iter::Restart.

/*
Return the number of nodes in the list which this list-iterator refers to.
*/
func (p *List_iter) ItemCount() int {
    //----------------------//
    // List_iter::ItemCount //
    //----------------------//
    if p == nil {
        return 0
    }
    if p.base == nil {
        return 0
    }
    return p.base.Length()
}   // End of function List_iter::ItemCount.

/*
Return the number of (nil, wrong, total) nodes in the list. See
List_base::ValidLength().
*/
func (p *List_iter) ItemCountValid() (int, int, int) {
    //------------------------------//
    //   List_iter::ItemCountValid  //
    //------------------------------//
    if p == nil {
        return 0, 0, 0
    }
    if p.base == nil {
        return 0, 0, 0
    }
    return p.base.ValidLength()
}   // End of function List_iter::ItemCountValid.

/*
List_iter::Next() returns the next node in the list. Checks are made to ensure
that the list is not corrupt. Any integrity errors cause the iteration to
terminate with a nil return value. For example, if the next node has been
removed from the list, and possibly appended to a different list,
List_iter::Next() will return a nil node-pointer and a non-nil error.

NOTE: The list should not be modified while iteration is occurring.
Results could be perplexing if the list is modified between Next-calls.
*/
func (p *List_iter) Next() (*List_node, error) {
    //----------------------//
    //    List_iter::Next   //
    //----------------------//
    /*------------------------------------------------------------------------------
      The field "current" initially indicates the outcome of the previous Next-call.
      If all goes well, this is incremented and returned to the caller.
      ------------------------------------------------------------------------------*/
    if p == nil {
        return nil, elist.New("List_base::Next: p == nil")
    }
    // If there's not list-base, there's nothing to do.
    if p.base == nil {
        return nil, elist.New("List_base::Next: p.base == nil")
    }
    if p.current == nil {
        p.current = p.base.first
        // Empty list.
        if p.current == nil {
            return nil, nil
        }
        // Corruption. The first node is not registered in a list!
        // Leave the current-pointer where it is to avoid infinite loops.
        if p.current.base == nil {
            return nil, elist.New("List_base::Next: p.current.base == nil")
        }
        // Corruption. The first node is in the wrong list!
        // Leave the current-pointer where it is to avoid infinite loops.
        if p.current.base != p.base {
            return nil, elist.New("List_base::Next: p.current.base != p.base")
        }
    } else {
        // The current node is not registered in a list!
        // Leave the current-pointer where it is to avoid infinite loops.
        if p.current.base == nil {
            return nil, elist.New("List_base::Next: p.current.base == nil")
        }
        // The current node is in the wrong list!
        // Leave the current-pointer where it is to avoid infinite loops.
        if p.current.base != p.base {
            return nil, elist.New("List_base::Next: p.current.base != p.base")
        }
        // End of the list.
        // Leave the current-pointer where it is to avoid infinite loops.
        if p.current.next == nil {
            return nil, nil
        }
        p.current = p.current.next
    }
    return p.current, nil
}   // End of function List_iter::Next.
