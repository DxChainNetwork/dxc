// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
pragma abicoder v2;

// OpenZeppelin Contracts (last updated v4.7.0) (utils/structs/EnumerableSet.sol)
/**

 * @dev Library for managing

 * https://en.wikipedia.org/wiki/Set_(abstract_data_type)[sets] of primitive

 * types.
   *

 * Sets have the following properties:
   *

 * - Elements are added, removed, and checked for existence in constant time

 * (O(1)).

 * - Elements are enumerated in O(n). No guarantees are made on the ordering.
     *

 * ```

   ```

 * contract Example {

 * // Add the library methods

 * using EnumerableSet for EnumerableSet.AddressSet;
    *

 * // Declare a set state variable

 * EnumerableSet.AddressSet private mySet;

 * }

 * ```
   *
   ```

 * As of v3.3.0, sets of type `bytes32` (`Bytes32Set`), `address` (`AddressSet`)

 * and `uint256` (`UintSet`) are supported.
   *

 * [WARNING]

 * ====

 * Trying to delete such a structure from storage will likely result in data corruption, rendering the structure unusable.

 * See https://github.com/ethereum/solidity/pull/11843[ethereum/solidity#11843] for more info.
   *

 * In order to clean an EnumerableSet, you can either remove all elements one by one or create a fresh instance using an array of EnumerableSet.

 * ====
   */
   library EnumerableSet {
   // To implement this library for multiple types with as little code
   // repetition as possible, we write it in terms of a generic Set type with
   // bytes32 values.
   // The Set implementation uses private functions, and user-facing
   // implementations (such as AddressSet) are just wrappers around the
   // underlying Set.
   // This means that we can only create new EnumerableSets for types that fit
   // in bytes32.

   struct Set {
       // Storage of set values
       bytes32[] _values;
       // Position of the value in the `values` array, plus 1 because index 0
       // means a value is not in the set.
       mapping(bytes32 => uint256) _indexes;
   }

   /**

    * @dev Add a value to a set. O(1).
      *
    * Returns true if the value was added to the set, that is if it was not
    * already present.
      */
      function _add(Set storage set, bytes32 value) private returns (bool) {
      if (!_contains(set, value)) {
          set._values.push(value);
          // The value is stored at length-1, but we add 1 to all indexes
          // and use 0 as a sentinel value
          set._indexes[value] = set._values.length;
          return true;
      } else {
          return false;
      }
      }

   /**

    * @dev Removes a value from a set. O(1).
      *

    * Returns true if the value was removed from the set, that is if it was

    * present.
      */
      function _remove(Set storage set, bytes32 value) private returns (bool) {
      // We read and store the value's index to prevent multiple reads from the same storage slot
      uint256 valueIndex = set._indexes[value];

      if (valueIndex != 0) {
          // Equivalent to contains(set, value)
          // To delete an element from the _values array in O(1), we swap the element to delete with the last one in
          // the array, and then remove the last element (sometimes called as 'swap and pop').
          // This modifies the order of the array, as noted in {at}.

          uint256 toDeleteIndex = valueIndex - 1;
          uint256 lastIndex = set._values.length - 1;

          if (lastIndex != toDeleteIndex) {
              bytes32 lastValue = set._values[lastIndex];

              // Move the last value to the index where the value to delete is
              set._values[toDeleteIndex] = lastValue;
              // Update the index for the moved value
              set._indexes[lastValue] = valueIndex; // Replace lastValue's index to valueIndex
          }

          // Delete the slot where the moved value was stored
          set._values.pop();

          // Delete the index for the deleted slot
          delete set._indexes[value];

          return true;

      } else {
          return false;
      }
      }

   /**

    * @dev Returns true if the value is in the set. O(1).
      */
      function _contains(Set storage set, bytes32 value)
      private
      view
      returns (bool)
      {
      return set._indexes[value] != 0;
      }

   /**

    * @dev Returns the number of values on the set. O(1).
      */
      function _length(Set storage set) private view returns (uint256) {
      return set._values.length;
      }

   /**

    * @dev Returns the value stored at position `index` in the set. O(1).
      *
    * Note that there are no guarantees on the ordering of values inside the
    * array, and it may change when more values are added or removed.
      *
    * Requirements:
      *
    * - `index` must be strictly less than {length}.
        */
        function _at(Set storage set, uint256 index)
        private
        view
        returns (bytes32)
        {
        return set._values[index];
        }

   /**

    * @dev Return the entire set in an array
      *
    * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
    * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
    * this function has an unbounded cost, and using it as part of a state-changing function may render the function
    * uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
      */
      function _values(Set storage set) private view returns (bytes32[] memory) {
      return set._values;
      }

   // Bytes32Set

   struct Bytes32Set {
       Set _inner;
   }

   /**

    * @dev Add a value to a set. O(1).
      *
    * Returns true if the value was added to the set, that is if it was not
    * already present.
      */
      function add(Bytes32Set storage set, bytes32 value)
      internal
      returns (bool)
      {
      return _add(set._inner, value);
      }

   /**

    * @dev Removes a value from a set. O(1).
      *
    * Returns true if the value was removed from the set, that is if it was
    * present.
      */
      function remove(Bytes32Set storage set, bytes32 value)
      internal
      returns (bool)
      {
      return _remove(set._inner, value);
      }

   /**

    * @dev Returns true if the value is in the set. O(1).
      */
      function contains(Bytes32Set storage set, bytes32 value)
      internal
      view
      returns (bool)
      {
      return _contains(set._inner, value);
      }

   /**

    * @dev Returns the number of values in the set. O(1).
      */
      function length(Bytes32Set storage set) internal view returns (uint256) {
      return _length(set._inner);
      }

   /**

    * @dev Returns the value stored at position `index` in the set. O(1).
      *
    * Note that there are no guarantees on the ordering of values inside the
    * array, and it may change when more values are added or removed.
      *
    * Requirements:
      *
    * - `index` must be strictly less than {length}.
        */
        function at(Bytes32Set storage set, uint256 index)
        internal
        view
        returns (bytes32)
        {
        return _at(set._inner, index);
        }

   /**

    * @dev Return the entire set in an array
      *
    * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed
    * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that
    * this function has an unbounded cost, and using it as part of a state-changing function may render the function
    * uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
      */
      function values(Bytes32Set storage set)
      internal
      view
      returns (bytes32[] memory)
      {
      return _values(set._inner);
      }

   // AddressSet

   struct AddressSet {
       Set _inner;
   }

   /**

    * @dev Add a value to a set. O(1).
      *
    * Returns true if the value was added to the set, that is if it was not
    * already present.
      */
      function add(AddressSet storage set, address value)
      internal
      returns (bool)
      {
      return _add(set._inner, bytes32(uint256(uint160(value))));
      }

   /**

    * @dev Removes a value from a set. O(1).
      *
    * Returns true if the value was removed from the set, that is if it was
    * present.
      */
      function remove(AddressSet storage set, address value)
      internal
      returns (bool)
      {
      return _remove(set._inner, bytes32(uint256(uint160(value))));
      }

   /**

    * @dev Returns true if the value is in the set. O(1).
      */
      function contains(AddressSet storage set, address value)
      internal
      view
      returns (bool)
      {
      return _contains(set._inner, bytes32(uint256(uint160(value))));
      }

   /**

    * @dev Returns the number of values in the set. O(1).
      */
      function length(AddressSet storage set) internal view returns (uint256) {
      return _length(set._inner);
      }

   /**

    * @dev Returns the value stored at position `index` in the set. O(1).
      *
    * Note that there are no guarantees on the ordering of values inside the
    * array, and it may change when more values are added or removed.
      *
    * Requirements:
      *
    * - `index` must be strictly less than {length}.
        */
        function at(AddressSet storage set, uint256 index)
        internal
        view
        returns (address)
        {
        return address(uint160(uint256(_at(set._inner, index))));
        }

   /**

    * @dev Return the entire set in an array
      *

    * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed

    * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that

    * this function has an unbounded cost, and using it as part of a state-changing function may render the function

    * uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
      */
      function values(AddressSet storage set)
      internal
      view
      returns (address[] memory)
      {
      bytes32[] memory store = _values(set._inner);
      address[] memory result;

      /// @solidity memory-safe-assembly
      assembly {
          result := store
      }

      return result;
      }

   // UintSet

   struct UintSet {
       Set _inner;
   }

   /**

    * @dev Add a value to a set. O(1).
      *
    * Returns true if the value was added to the set, that is if it was not
    * already present.
      */
      function add(UintSet storage set, uint256 value) internal returns (bool) {
      return _add(set._inner, bytes32(value));
      }

   /**

    * @dev Removes a value from a set. O(1).
      *
    * Returns true if the value was removed from the set, that is if it was
    * present.
      */
      function remove(UintSet storage set, uint256 value)
      internal
      returns (bool)
      {
      return _remove(set._inner, bytes32(value));
      }

   /**

    * @dev Returns true if the value is in the set. O(1).
      */
      function contains(UintSet storage set, uint256 value)
      internal
      view
      returns (bool)
      {
      return _contains(set._inner, bytes32(value));
      }

   /**

    * @dev Returns the number of values on the set. O(1).
      */
      function length(UintSet storage set) internal view returns (uint256) {
      return _length(set._inner);
      }

   /**

    * @dev Returns the value stored at position `index` in the set. O(1).
      *
    * Note that there are no guarantees on the ordering of values inside the
    * array, and it may change when more values are added or removed.
      *
    * Requirements:
      *
    * - `index` must be strictly less than {length}.
        */
        function at(UintSet storage set, uint256 index)
        internal
        view
        returns (uint256)
        {
        return uint256(_at(set._inner, index));
        }

   /**

    * @dev Return the entire set in an array
      *

    * WARNING: This operation will copy the entire storage to memory, which can be quite expensive. This is designed

    * to mostly be used by view accessors that are queried without any gas fees. Developers should keep in mind that

    * this function has an unbounded cost, and using it as part of a state-changing function may render the function

    * uncallable if the set grows to a point where copying to memory consumes too much gas to fit in a block.
      */
      function values(UintSet storage set)
      internal
      view
      returns (uint256[] memory)
      {
      bytes32[] memory store = _values(set._inner);
      uint256[] memory result;

      /// @solidity memory-safe-assembly
      assembly {
          result := store
      }

      return result;
      }
      }

// OpenZeppelin Contracts (last updated v4.7.0) (utils/Address.sol)
/**

 * @dev Collection of functions related to the address type
   */
   library Address {
   /**

    * @dev Returns true if `account` is a contract.
      *

    * [IMPORTANT]

    * ====

    * It is unsafe to assume that an address for which this function returns

    * false is an externally-owned account (EOA) and not a contract.
      *

    * Among others, `isContract` will return false for the following

    * types of addresses:
      *

    * - an externally-owned account

    * - a contract in construction

    * - an address where a contract will be created

    * - an address where a contract lived, but was destroyed

    * ====
      *

    * [IMPORTANT]

    * ====

    * You shouldn't rely on `isContract` to protect against flash loan attacks!
      *

    * Preventing calls from contracts is highly discouraged. It breaks composability, breaks support for smart wallets

    * like Gnosis Safe, and does not provide security since it can be circumvented by calling from a contract

    * constructor.

    * ====
      */
      function isContract(address account) internal view returns (bool) {
      // This method relies on extcodesize/address.code.length, which returns 0
      // for contracts in construction, since the code is only stored at the end
      // of the constructor execution.

      return account.code.length > 0;
      }

   /**

    * @dev Replacement for Solidity's `transfer`: sends `amount` wei to

    * `recipient`, forwarding all available gas and reverting on errors.
      *

    * https://eips.ethereum.org/EIPS/eip-1884[EIP1884] increases the gas cost

    * of certain opcodes, possibly making contracts go over the 2300 gas limit

    * imposed by `transfer`, making them unable to receive funds via

    * `transfer`. {sendValue} removes this limitation.
      *

    * https://diligence.consensys.net/posts/2019/09/stop-using-soliditys-transfer-now/[Learn more].
      *

    * IMPORTANT: because control is transferred to `recipient`, care must be

    * taken to not create reentrancy vulnerabilities. Consider using

    * {ReentrancyGuard} or the

    * https://solidity.readthedocs.io/en/v0.5.11/security-considerations.html#use-the-checks-effects-interactions-pattern[checks-effects-interactions pattern].
      */
      function sendValue(address payable recipient, uint256 amount) internal {
      require(
          address(this).balance >= amount,
          "Address: insufficient balance"
      );

      (bool success, ) = recipient.call{value: amount}("");
      require(
          success,
          "Address: unable to send value, recipient may have reverted"
      );
      }

   /**

    * @dev Performs a Solidity function call using a low level `call`. A
    * plain `call` is an unsafe replacement for a function call: use this
    * function instead.
      *
    * If `target` reverts with a revert reason, it is bubbled up by this
    * function (like regular Solidity function calls).
      *
    * Returns the raw returned data. To convert to the expected return value,
    * use https://solidity.readthedocs.io/en/latest/units-and-global-variables.html?highlight=abi.decode#abi-encoding-and-decoding-functions[`abi.decode`].
      *
    * Requirements:
      *
    * - `target` must be a contract.
    * - calling `target` with `data` must not revert.
        *
    * _Available since v3.1._
      */
      function functionCall(address target, bytes memory data)
      internal
      returns (bytes memory)
      {
      return functionCall(target, data, "Address: low-level call failed");
      }

   /**

    * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`], but with
    * `errorMessage` as a fallback revert reason when `target` reverts.
      *
    * _Available since v3.1._
      */
      function functionCall(
      address target,
      bytes memory data,
      string memory errorMessage
      ) internal returns (bytes memory) {
      return functionCallWithValue(target, data, 0, errorMessage);
      }

   /**

    * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`],
    * but also transferring `value` wei to `target`.
      *
    * Requirements:
      *
    * - the calling contract must have an ETH balance of at least `value`.
    * - the called Solidity function must be `payable`.
        *
    * _Available since v3.1._
      */
      function functionCallWithValue(
      address target,
      bytes memory data,
      uint256 value
      ) internal returns (bytes memory) {
      return
          functionCallWithValue(
              target,
              data,
              value,
              "Address: low-level call with value failed"
          );
      }

   /**

    * @dev Same as {xref-Address-functionCallWithValue-address-bytes-uint256-}[`functionCallWithValue`], but

    * with `errorMessage` as a fallback revert reason when `target` reverts.
      *

    * _Available since v3.1._
      */
      function functionCallWithValue(
      address target,
      bytes memory data,
      uint256 value,
      string memory errorMessage
      ) internal returns (bytes memory) {
      require(
          address(this).balance >= value,
          "Address: insufficient balance for call"
      );
      require(isContract(target), "Address: call to non-contract");

      (bool success, bytes memory returndata) = target.call{value: value}(
          data
      );
      return verifyCallResult(success, returndata, errorMessage);
      }

   /**

    * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`],
    * but performing a static call.
      *
    * _Available since v3.3._
      */
      function functionStaticCall(address target, bytes memory data)
      internal
      view
      returns (bytes memory)
      {
      return
          functionStaticCall(
              target,
              data,
              "Address: low-level static call failed"
          );
      }

   /**

    * @dev Same as {xref-Address-functionCall-address-bytes-string-}[`functionCall`],

    * but performing a static call.
      *

    * _Available since v3.3._
      */
      function functionStaticCall(
      address target,
      bytes memory data,
      string memory errorMessage
      ) internal view returns (bytes memory) {
      require(isContract(target), "Address: static call to non-contract");

      (bool success, bytes memory returndata) = target.staticcall(data);
      return verifyCallResult(success, returndata, errorMessage);
      }

   /**

    * @dev Same as {xref-Address-functionCall-address-bytes-}[`functionCall`],
    * but performing a delegate call.
      *
    * _Available since v3.4._
      */
      function functionDelegateCall(address target, bytes memory data)
      internal
      returns (bytes memory)
      {
      return
          functionDelegateCall(
              target,
              data,
              "Address: low-level delegate call failed"
          );
      }

   /**

    * @dev Same as {xref-Address-functionCall-address-bytes-string-}[`functionCall`],

    * but performing a delegate call.
      *

    * _Available since v3.4._
      */
      function functionDelegateCall(
      address target,
      bytes memory data,
      string memory errorMessage
      ) internal returns (bytes memory) {
      require(isContract(target), "Address: delegate call to non-contract");

      (bool success, bytes memory returndata) = target.delegatecall(data);
      return verifyCallResult(success, returndata, errorMessage);
      }

   /**

    * @dev Tool to verifies that a low level call was successful, and revert if it wasn't, either by bubbling the
    * revert reason using the provided one.
      *
    * _Available since v4.3._
      */
      function verifyCallResult(
      bool success,
      bytes memory returndata,
      string memory errorMessage
      ) internal pure returns (bytes memory) {
      if (success) {
          return returndata;
      } else {
          // Look for revert reason and bubble it up if present
          if (returndata.length > 0) {
              // The easiest way to bubble the revert reason is using memory via assembly
              /// @solidity memory-safe-assembly
              assembly {
                  let returndata_size := mload(returndata)
                  revert(add(32, returndata), returndata_size)
              }
          } else {
              revert(errorMessage);
          }
      }
      }
      }

// OpenZeppelin Contracts (last updated v4.7.0) (proxy/utils/Initializable.sol)
/**

 * @dev This is a base contract to aid in writing upgradeable contracts, or any kind of contract that will be deployed

 * behind a proxy. Since proxied contracts do not make use of a constructor, it's common to move constructor logic to an

 * external initializer function, usually called `initialize`. It then becomes necessary to protect this initializer

 * function so it can only be called once. The {initializer} modifier provided by this contract will have this effect.
   *

 * The initialization functions use a version number. Once a version number is used, it is consumed and cannot be

 * reused. This mechanism prevents re-execution of each "step" but allows the creation of new initialization steps in

 * case an upgrade adds a module that needs to be initialized.
   *

 * For example:
   *

 * [.hljs-theme-light.nopadding]

 * ```

   ```

 * contract MyToken is ERC20Upgradeable {

 * function initialize() initializer public {

 * __ERC20_init("MyToken", "MTK");

 * }

 * }

 * contract MyTokenV2 is MyToken, ERC20PermitUpgradeable {

 * function initializeV2() reinitializer(2) public {

 * __ERC20Permit_init("MyToken");

 * }

 * }

 * ```
   *
   ```

 * TIP: To avoid leaving the proxy in an uninitialized state, the initializer function should be called as early as

 * possible by providing the encoded function call as the `_data` argument to {ERC1967Proxy-constructor}.
   *

 * CAUTION: When used with inheritance, manual care must be taken to not invoke a parent initializer twice, or to ensure

 * that all initializers are idempotent. This is not verified automatically as constructors are by Solidity.
   *

 * [CAUTION]

 * ====

 * Avoid leaving a contract uninitialized.
   *

 * An uninitialized contract can be taken over by an attacker. This applies to both a proxy and its implementation

 * contract, which may impact the proxy. To prevent the implementation contract from being used, you should invoke

 * the {_disableInitializers} function in the constructor to automatically lock it when it is deployed:
   *

 * [.hljs-theme-light.nopadding]

 * ```

   ```

 * /// @custom:oz-upgrades-unsafe-allow constructor

 * constructor() {

 * _disableInitializers();

 * }

 * ```

   ```

 * ====
   */
   abstract contract Initializable {
   /**

    * @dev Indicates that the contract has been initialized.
    * @custom:oz-retyped-from bool
      */
      uint8 private _initialized;

   /**

    * @dev Indicates that the contract is in the process of being initialized.
      */
      bool private _initializing;

   /**

    * @dev Triggered when the contract has been initialized or reinitialized.
      */
      event Initialized(uint8 version);

   /**

    * @dev A modifier that defines a protected initializer function that can be invoked at most once. In its scope,
    * `onlyInitializing` functions can be used to initialize parent contracts. Equivalent to `reinitializer(1)`.
      */
      modifier initializer() {
      bool isTopLevelCall = !_initializing;
      require(
          (isTopLevelCall && _initialized < 1) ||
              (!Address.isContract(address(this)) && _initialized == 1),
          "Initializable: contract is already initialized"
      );
      _initialized = 1;
      if (isTopLevelCall) {
          _initializing = true;
      }
      _;
      if (isTopLevelCall) {
          _initializing = false;
          emit Initialized(1);
      }
      }

   /**

    * @dev A modifier that defines a protected reinitializer function that can be invoked at most once, and only if the
    * contract hasn't been initialized to a greater version before. In its scope, `onlyInitializing` functions can be
    * used to initialize parent contracts.
      *
    * `initializer` is equivalent to `reinitializer(1)`, so a reinitializer may be used after the original
    * initialization step. This is essential to configure modules that are added through upgrades and that require
    * initialization.
      *
    * Note that versions can jump in increments greater than 1; this implies that if multiple reinitializers coexist in
    * a contract, executing them in the right order is up to the developer or operator.
      */
      modifier reinitializer(uint8 version) {
      require(
          !_initializing && _initialized < version,
          "Initializable: contract is already initialized"
      );
      _initialized = version;
      _initializing = true;
      _;
      _initializing = false;
      emit Initialized(version);
      }

   /**

    * @dev Modifier to protect an initialization function so that it can only be invoked by functions with the
    * {initializer} and {reinitializer} modifiers, directly or indirectly.
      */
      modifier onlyInitializing() {
      require(_initializing, "Initializable: contract is not initializing");
      _;
      }

   /**

    * @dev Locks the contract, preventing any future reinitialization. This cannot be part of an initializer call.
    * Calling this in the constructor of a contract will prevent that contract from being initialized or reinitialized
    * to any version. It is recommended to use this to lock implementation contracts that are designed to be called
    * through proxies.
      */
      function _disableInitializers() internal virtual {
      require(!_initializing, "Initializable: contract is initializing");
      if (_initialized < type(uint8).max) {
          _initialized = type(uint8).max;
          emit Initialized(type(uint8).max);
      }
      }
      }

// OpenZeppelin Contracts v4.4.1 (security/ReentrancyGuard.sol)
/**

 * @dev Contract module that helps prevent reentrant calls to a function.
   *

 * Inheriting from `ReentrancyGuard` will make the {nonReentrant} modifier

 * available, which can be applied to functions to make sure there are no nested

 * (reentrant) calls to them.
   *

 * Note that because there is a single `nonReentrant` guard, functions marked as

 * `nonReentrant` may not call one another. This can be worked around by making

 * those functions `private`, and then adding `external` `nonReentrant` entry

 * points to them.
   *

 * TIP: If you would like to learn more about reentrancy and alternative ways

 * to protect against it, check out our blog post

 * https://blog.openzeppelin.com/reentrancy-after-istanbul/[Reentrancy After Istanbul].
   */
   abstract contract ReentrancyGuard {
   // Booleans are more expensive than uint256 or any type that takes up a full
   // word because each write operation emits an extra SLOAD to first read the
   // slot's contents, replace the bits taken up by the boolean, and then write
   // back. This is the compiler's defense against contract upgrades and
   // pointer aliasing, and it cannot be disabled.

   // The values being non-zero value makes deployment a bit more expensive,
   // but in exchange the refund on every call to nonReentrant will be lower in
   // amount. Since refunds are capped to a percentage of the total
   // transaction's gas, it is best to keep them low in cases like this one, to
   // increase the likelihood of the full refund coming into effect.
   uint256 private constant _NOT_ENTERED = 1;
   uint256 private constant _ENTERED = 2;

   uint256 private _status;

   constructor() {
       _status = _NOT_ENTERED;
   }

   /**

    * @dev Prevents a contract from calling itself, directly or indirectly.

    * Calling a `nonReentrant` function from another `nonReentrant`

    * function is not supported. It is possible to prevent this from happening

    * by making the `nonReentrant` function external, and making it call a

    * `private` function that does the actual work.
      */
      modifier nonReentrant() {
      // On the first call to nonReentrant, _notEntered will be true
      require(_status != _ENTERED, "ReentrancyGuard: reentrant call");

      // Any calls to nonReentrant after this point will fail
      _status = _ENTERED;

      _;

      // By storing the original value once again, a refund is triggered (see
      // https://eips.ethereum.org/EIPS/eip-2200)
      _status = _NOT_ENTERED;
      }
      }

contract Base {
    uint256 public constant BLOCK_SECONDS = 6;
    /// @notice min rate. base on 100
    uint8 public constant MIN_RATE = 70;
    /// @notice max rate. base on 100
    uint8 public constant MAX_RATE = 100;

    /// @notice 10 * 60 / BLOCK_SECONDS
    uint256 public constant EPOCH_BLOCKS = 14400;
    /// @notice min deposit for validator
    uint256 public constant MIN_DEPOSIT = 4e7 ether;
    uint256 public constant MAX_PUNISH_COUNT = 139;

    /// @notice use blocks as units in code: RATE_SET_LOCK_EPOCHS * EPOCH_BLOCKS
    uint256 public constant RATE_SET_LOCK_EPOCHS = 1;
    /// @notice use blocks as units in code: VALIDATOR_UNSTAKE_LOCK_EPOCHS * EPOCH_BLOCKS
    uint256 public constant VALIDATOR_UNSTAKE_LOCK_EPOCHS = 1;
    /// @notice use blocks as units in code: PROPOSAL_DURATION_EPOCHS * EPOCH_BLOCKS
    uint256 public constant PROPOSAL_DURATION_EPOCHS = 7;
    /// @notice use epoch as units in code: VALIDATOR_REWARD_LOCK_EPOCHS
    uint256 public constant VALIDATOR_REWARD_LOCK_EPOCHS = 7;
    /// @notice use epoch as units in code: VOTE_CANCEL_EPOCHS
    uint256 public constant VOTE_CANCEL_EPOCHS = 1;

    uint256 public constant MAX_VALIDATORS_COUNT = 210;
    uint256 public constant MAX_VALIDATOR_DETAIL_LENGTH = 1000;
    uint256 public constant MAX_VALIDATOR_NAME_LENGTH = 100;

    // total deposit
    uint256 public constant TOTAL_DEPOSIT_LV1 = 1e18 * 1e8 * 150;
    uint256 public constant TOTAL_DEPOSIT_LV2 = 1e18 * 1e8 * 200;
    uint256 public constant TOTAL_DEPOSIT_LV3 = 1e18 * 1e8 * 250;
    uint256 public constant TOTAL_DEPOSIT_LV4 = 1e18 * 1e8 * 300;
    uint256 public constant TOTAL_DEPOSIT_LV5 = 1e18 * 1e8 * 350;

    // block reward
    uint256 public constant REWARD_DEPOSIT_UNDER_LV1 = 1e15 * 95250;
    uint256 public constant REWARD_DEPOSIT_FROM_LV1_TO_LV2 = 1e15 * 128250;
    uint256 public constant REWARD_DEPOSIT_FROM_LV2_TO_LV3 = 1e15 * 157125;
    uint256 public constant REWARD_DEPOSIT_FROM_LV3_TO_LV4 = 1e15 * 180750;
    uint256 public constant REWARD_DEPOSIT_FROM_LV4_TO_LV5 = 1e15 * 199875;
    uint256 public constant REWARD_DEPOSIT_OVER_LV5 = 1e15 * 214125;

    // validator count
    uint256 public constant MAX_VALIDATOR_COUNT_LV1 = 21;
    uint256 public constant MAX_VALIDATOR_COUNT_LV2 = 33;
    uint256 public constant MAX_VALIDATOR_COUNT_LV3 = 66;
    uint256 public constant MAX_VALIDATOR_COUNT_LV4 = 99;
    uint256 public constant MIN_LEVEL_VALIDATOR_COUNT = 60;
    uint256 public constant MEDIUM_LEVEL_VALIDATOR_COUNT = 90;
    uint256 public constant MAX_LEVEL_VALIDATOR_COUNT = 120;

    // dead address
    address public constant BLACK_HOLE_ADDRESS =
        0x0000000000000000000000000000000000000000;

    uint256 public constant SAFE_MULTIPLIER = 1e18;

    modifier onlySystem() {
        require(tx.gasprice == 0, "Prohibit external calls");
        _;
    }

    modifier onlyMiner() {
        require(msg.sender == block.coinbase, "msg.sender error");
        _;
    }

    /**
     * @dev return current epoch
     */
    function currentEpoch() public view returns (uint256) {
        return block.number / EPOCH_BLOCKS;
    }

}

interface ISystemRewards {
    function epochs(uint256 _epoch)
        external
        view
        returns (
            uint256,
            uint256,
            uint256,
            uint256
        );

    function updateValidatorWhileElect(
        address _val,
        uint8 _rate,
        uint256 _newEpoch
    ) external;

    function updateEpochWhileElect(
        uint256 _tvl,
        uint256 _valCount,
        uint256 _effictiveValCount,
        uint256 _newEpoch
    ) external;

    function updateValidatorWhileEpochEnd(address _val, uint256 _votes)
        external;

    function getRewardPerVote(address _val) external view returns (uint256);

}

interface IProposals {}

interface INodeVote {
    function totalVotes() external view returns (uint256);
}

contract Validators is Base, Initializable, ReentrancyGuard {
    using EnumerableSet for EnumerableSet.AddressSet;
    using Address for address;

    enum ValidatorStatus {
        canceled,
        canceling,
        cancelQueue,
        kickout,
        effictive
    }

    address[] public curEpochValidators;
    mapping(address => uint256) public curEpochValidatorsIdMap;

    EnumerableSet.AddressSet effictiveValidators;

    /// @notice canceled、canceling、kickout
    EnumerableSet.AddressSet invalidValidators;

    struct Validator {
        ValidatorStatus status;
        uint256 deposit;
        uint8 rate;
        /// @notice name
        string name;
        /// @notice details
        string details;
        uint256 votes;
        uint256 unstakeLockingEndBlock;
        uint256 rateSettLockingEndBlock;
    }

    mapping(address => Validator) _validators;

    mapping(address => EnumerableSet.AddressSet) validatorToVoters;

    /// @notice all validators deposit
    uint256 public totalDeposit;

    /// @notice  cancel validator queue
    EnumerableSet.AddressSet cancelingValidators;

    /// @notice SystemRewards contract
    ISystemRewards public sysRewards;

    /// @notice Proposals contract
    IProposals public proposals;

    /// @notice NodeVote contract
    INodeVote public nodeVote;

    event LogAddValidator(
        address indexed _val,
        uint256 _deposit,
        uint256 _rate
    );
    event LogUpdateValidatorDeposit(address indexed _val, uint256 _deposit);
    event LogUpdateValidatorRate(
        address indexed _val,
        uint8 _preRate,
        uint8 _rate
    );
    event LogUnstakeValidator(address indexed _val);
    event LogRedeemValidator(address indexed _val);
    event LogRestoreValidator(address indexed _val);

    /**
     * @dev only Proposals contract address
     */
    modifier onlyProposalsC() {
        require(
            msg.sender == address(proposals),
            "Validators: not Proposals contract address"
        );
        _;
    }

    /**
     * @dev only SystemRewards contract address
     */
    modifier onlySysRewardsC() {
        require(
            msg.sender == address(sysRewards),
            "Validators: not SystemRewards contract address"
        );
        _;
    }

    /**
     * @dev only NodeVote contract address
     */
    modifier onlyNodeVoteC() {
        require(
            msg.sender == address(nodeVote),
            "Validators: not NodeVote contract address"
        );
        _;
    }

    function initialize(
        address _proposal,
        address _sysReward,
        address _nodeVote,
        address _initVal,
        uint256 _initDeposit,
        uint8 _initRate,
        string memory _name,
        string memory _details
    ) external payable onlySystem initializer {
        sysRewards = ISystemRewards(_sysReward);
        proposals = IProposals(_proposal);
        nodeVote = INodeVote(_nodeVote);

        require(!_initVal.isContract(), "Validators: validator address error");
        require(
            msg.value == _initDeposit && _initDeposit >= MIN_DEPOSIT,
            "Validators: deposit or value error"
        );
        require(
            _initRate >= MIN_RATE && _initRate <= MAX_RATE,
            "Validators: Rate must greater than MIN_RATE and less than MAX_RATE"
        );

        Validator storage val = _validators[_initVal];
        val.status = ValidatorStatus.effictive;
        val.deposit = _initDeposit;
        val.rate = _initRate;
        val.name = _name;
        val.details = _details;

        effictiveValidators.add(_initVal);
        totalDeposit += _initDeposit;

        curEpochValidators.push(_initVal);
        curEpochValidatorsIdMap[_initVal] = curEpochValidators.length;

        uint256 curEpoch = currentEpoch();
        sysRewards.updateValidatorWhileElect(_initVal, _initRate, curEpoch);
        sysRewards.updateEpochWhileElect(
            totalDeposit,
            curEpochValidators.length,
            effictiveValidators.length(),
            curEpoch
        );

        emit LogAddValidator(_initVal, _initDeposit, _initRate);
    }

    /**
     * @dev get voter
     */
    function getValidatorVoters(
        address _val,
        uint256 page,
        uint256 size
    ) external view returns (address[] memory) {
        require(page > 0 && size > 0, "Validators: Requests param error");
        EnumerableSet.AddressSet storage voters = validatorToVoters[_val];
        uint256 start = (page - 1) * size;
        if (voters.length() < start) {
            size = 0;
        } else {
            uint256 length = voters.length() - start;
            if (length < size) {
                size = length;
            }
        }

        address[] memory vals = new address[](size);
        for (uint256 i = 0; i < size; i++) {
            vals[i] = voters.at(i + start);
        }
        return vals;
    }

    /**
     * @dev return voters count of validator
     */
    function validatorVotersLength(address _val) public view returns (uint256) {
        return validatorToVoters[_val].length();
    }

    /**
     * @dev return validator info
     */
    function validators(address _val) external view returns (Validator memory) {
        return _validators[_val];
    }

    /**
     * @dev batch query validator info
     */
    function batchValidators(address[] memory _vals)
        external
        view
        returns (Validator[] memory)
    {
        uint256 len = _vals.length;
        Validator[] memory valInfos = new Validator[](len);

        for (uint256 i = 0; i < len; i++) {
            valInfos[i] = _validators[_vals[i]];
        }
        return valInfos;
    }

    /**
     * @dev return curEpochValidators
     */
    function getCurEpochValidators() external view returns (address[] memory) {
        return curEpochValidators;
    }

    /**
     * @dev True: effictive
     */
    function isEffictiveValidator(address addr) external view returns (bool) {
        return _validators[addr].status == ValidatorStatus.effictive;
    }

    /**
     * @dev return effictive Validators count
     */
    function effictiveValsLength() public view returns (uint256) {
        return effictiveValidators.length();
    }

    /**
     * @dev return all effictive Validators
     */
    function getEffictiveValidators() public view returns (address[] memory) {
        uint256 len = effictiveValidators.length();
        address[] memory vals = new address[](len);

        for (uint256 i = 0; i < len; i++) {
            vals[i] = effictiveValidators.at(i);
        }
        return vals;
    }

    function getEffictiveValidatorsWithPage(uint256 page, uint256 size)
        public
        view
        returns (address[] memory)
    {
        require(page > 0 && size > 0, "Validators: Requests param error");
        uint256 len = effictiveValidators.length();
        uint256 start = (page - 1) * size;
        if (len < start) {
            size = 0;
        } else {
            uint256 length = len - start;
            if (length < size) {
                size = length;
            }
        }

        address[] memory vals = new address[](size);
        for (uint256 i = 0; i < size; i++) {
            vals[i] = effictiveValidators.at(i + start);
        }
        return vals;
    }

    /**
     * @dev return invalid Validators count
     */
    function invalidValsLength() public view returns (uint256) {
        return invalidValidators.length();
    }

    /**
     * @dev return all invalid Validators
     */
    function getInvalidValidators() public view returns (address[] memory) {
        uint256 len = invalidValidators.length();
        address[] memory vals = new address[](len);

        for (uint256 i = 0; i < len; i++) {
            vals[i] = invalidValidators.at(i);
        }
        return vals;
    }

    function getInvalidValidatorsWithPage(uint256 page, uint256 size)
        public
        view
        returns (address[] memory)
    {
        require(page > 0 && size > 0, "Validators: Requests param error");
        uint256 len = invalidValidators.length();
        uint256 start = (page - 1) * size;
        if (len < start) {
            size = 0;
        } else {
            uint256 length = len - start;
            if (length < size) {
                size = length;
            }
        }

        address[] memory vals = new address[](size);
        for (uint256 i = 0; i < size; i++) {
            vals[i] = invalidValidators.at(i + start);
        }
        return vals;
    }

    /**
     * @dev return canceling validators count
     */
    function CancelQueueValidatorsLength() public view returns (uint256) {
        return cancelingValidators.length();
    }

    /**
     * @dev return Cancel Queue Validators
     */
    function getCancelQueueValidators() public view returns (address[] memory) {
        uint256 len = cancelingValidators.length();
        address[] memory vals = new address[](len);

        for (uint256 i = 0; i < len; i++) {
            vals[i] = cancelingValidators.at(i);
        }
        return vals;
    }

    /**
     * @dev update validator deposit
     */
    function updateValidatorDeposit(uint256 _deposit)
        external
        payable
        nonReentrant
    {
        Validator storage val = _validators[msg.sender];
        require(
            val.status == ValidatorStatus.effictive,
            "Validators: illegal msg.sender"
        );
        if (_deposit >= val.deposit) {
            require(
                msg.value >= _deposit - val.deposit,
                "Validators: illegal deposit"
            );
            uint256 sub = _deposit - val.deposit;
            totalDeposit += sub;
            val.deposit = _deposit;
            payable(msg.sender).transfer(msg.value - sub);
        } else {
            require(_deposit >= MIN_DEPOSIT, "Validators: illegal deposit");
            uint256 sub = val.deposit - _deposit;
            payable(msg.sender).transfer(sub);
            val.deposit = _deposit;
            totalDeposit -= sub;
        }

        emit LogUpdateValidatorDeposit(msg.sender, val.deposit);
    }

    /**
     * @dev update validator rate
     */
    function updateValidatorRate(uint8 _rate) external nonReentrant {
        Validator storage val = _validators[msg.sender];
        require(
            val.status == ValidatorStatus.effictive,
            "Validators: illegal msg.sender"
        );
        require(
            val.rateSettLockingEndBlock < block.number,
            "Validators: illegal rate set block"
        );
        require(
            _rate >= MIN_RATE && val.rate <= MAX_RATE,
            "Validators: illegal Allocation ratio"
        );
        uint8 preRate = val.rate;
        val.rate = _rate;
        val.rateSettLockingEndBlock =
            block.number +
            RATE_SET_LOCK_EPOCHS *
            EPOCH_BLOCKS;

        emit LogUpdateValidatorRate(msg.sender, preRate, _rate);
    }

    /**
     * @dev update validator name and details
     */
    function updateValidatorNameDetails(
        string memory _name,
        string memory _details
    ) external nonReentrant {
        Validator storage val = _validators[msg.sender];
        require(
            bytes(_details).length <= MAX_VALIDATOR_DETAIL_LENGTH,
            "Validators: Details is too long"
        );
        require(
            bytes(_name).length <= MAX_VALIDATOR_NAME_LENGTH,
            "Validators: name is too long"
        );
        val.name = _name;
        val.details = _details;
    }

    function addValidatorFromProposal(
        address _val,
        uint256 _deposit,
        uint8 _rate,
        string memory _name,
        string memory _details
    ) external payable onlyProposalsC {
        require(!_val.isContract(), "Validators: validator address error");
        require(
            msg.value == _deposit,
            "Validators: deposit not equal msg.value"
        );

        Validator storage val = _validators[_val];
        require(
            val.status == ValidatorStatus.canceled,
            "Validators: validator status error"
        );

        val.status = ValidatorStatus.effictive;
        val.deposit = _deposit;
        val.rate = _rate;
        val.name = _name;
        val.details = _details;

        effictiveValidators.add(_val);
        invalidValidators.remove(_val);
        totalDeposit += _deposit;

        emit LogAddValidator(_val, _deposit, _rate);
    }

    function kickoutValidator(address _val) external onlySysRewardsC {
        Validator storage val = _validators[_val];
        require(
            val.status == ValidatorStatus.effictive ||
                val.status == ValidatorStatus.kickout,
            "Validators: validator status error"
        );
        val.status = ValidatorStatus.kickout;
        if (effictiveValidators.contains(_val)) {
            effictiveValidators.remove(_val);
            invalidValidators.add(_val);
            totalDeposit -= val.deposit;
        }
    }

    function restore() external nonReentrant {
        require(
            effictiveValidators.length() < MAX_VALIDATORS_COUNT,
            "Validators: length of the validator must be less than MAX_VALIDATORS_COUNT"
        );
        Validator storage val = _validators[msg.sender];
        require(
            !cancelingValidators.contains(msg.sender),
            "Validators: this validator is canceling"
        );
        require(
            val.status == ValidatorStatus.kickout,
            "Validators: validator must be kickout"
        );
        val.status = ValidatorStatus.effictive;
        effictiveValidators.add(msg.sender);
        invalidValidators.remove(msg.sender);
        totalDeposit += val.deposit;

        emit LogRestoreValidator(msg.sender);
    }

    function unstake() external nonReentrant {
        Validator storage val = _validators[msg.sender];
        require(
            val.status == ValidatorStatus.effictive ||
                val.status == ValidatorStatus.kickout,
            "Validators: illegal msg.sender"
        );
        if (curEpochValidatorsIdMap[msg.sender] == 0) {
            cancelingValidators.remove(msg.sender);
            val.status = ValidatorStatus.canceling;
            val.unstakeLockingEndBlock =
                block.number +
                VALIDATOR_UNSTAKE_LOCK_EPOCHS *
                EPOCH_BLOCKS;

            if (effictiveValidators.contains(msg.sender)) {
                effictiveValidators.remove(msg.sender);
                invalidValidators.add(msg.sender);
                totalDeposit -= val.deposit;
            }
        } else {
            val.status = ValidatorStatus.cancelQueue;
            cancelingValidators.add(msg.sender);
        }
        emit LogUnstakeValidator(msg.sender);
    }

    function _cancelValidatorWhileElect() internal {
        for (uint256 i = 0; i < cancelingValidators.length(); i++) {
            address _val = cancelingValidators.at(0);

            Validator storage val = _validators[_val];
            val.status = ValidatorStatus.canceling;
            val.unstakeLockingEndBlock =
                block.number +
                VALIDATOR_UNSTAKE_LOCK_EPOCHS *
                EPOCH_BLOCKS;

            cancelingValidators.remove(_val);

            if (effictiveValidators.contains(_val)) {
                effictiveValidators.remove(_val);
                invalidValidators.add(_val);
                totalDeposit -= val.deposit;
            }
        }
    }

    function redeem() external nonReentrant {
        Validator storage val = _validators[msg.sender];
        require(
            val.unstakeLockingEndBlock < block.number,
            "Validators: illegal redeem block"
        );
        require(
            val.status == ValidatorStatus.canceling &&
                curEpochValidatorsIdMap[msg.sender] == 0,
            "Validators: illegal msg.sender"
        );

        val.status = ValidatorStatus.canceled;
        payable(msg.sender).transfer(val.deposit);
        val.deposit = 0;
        val.unstakeLockingEndBlock = 0;
        val.rateSettLockingEndBlock = 0;
        invalidValidators.remove(msg.sender);

        emit LogRedeemValidator(msg.sender);
    }

    function voteValidator(
        address _voter,
        address _val,
        uint256 _votes
    ) external payable onlyNodeVoteC {
        _validators[_val].votes += _votes;
        validatorToVoters[_val].add(_voter);
    }

    function cancelVoteValidator(
        address _voter,
        address _val,
        uint256 _votes,
        bool _clear
    ) external onlyNodeVoteC {
        _validators[_val].votes -= _votes;
        if (_clear) {
            validatorToVoters[_val].remove(_voter);
        }
    }

    function tryElect() external onlySysRewardsC {
        _cancelValidatorWhileElect();

        uint256 nextEpochValCount = nextEpochValidatorCount();
        uint256 effictiveLen = effictiveValidators.length();

        for (uint256 i = 0; i < curEpochValidators.length; i++) {
            address _val = curEpochValidators[i];
            sysRewards.updateValidatorWhileEpochEnd(
                _val,
                _validators[_val].votes
            );
            delete curEpochValidatorsIdMap[_val];
        }
        delete curEpochValidators;

        uint256 total = 0;
        for (uint256 i = 0; i < effictiveLen; i++) {
            address val = effictiveValidators.at(i);
            total += _validators[val].votes + _validators[val].deposit;
        }

        uint256 totalTemp = total;
        uint256 nextEpoch = currentEpoch() + 1;

        if (nextEpochValCount >= effictiveLen) {
            for (uint256 i = 0; i < effictiveLen; i++) {
                address val = effictiveValidators.at(i);
                curEpochValidators.push(val);
                curEpochValidatorsIdMap[val] = curEpochValidators.length;
                sysRewards.updateValidatorWhileElect(
                    val,
                    _validators[val].rate,
                    nextEpoch
                );
            }
        } else {
            // for-loop tryElect
            for (uint256 i = 0; i < nextEpochValCount; i++) {
                if (total <= 0) break;
                // get random number
                uint256 randDeposit = rand(total, i);

                for (uint256 j = 0; j < effictiveLen; j++) {
                    address val = effictiveValidators.at(j);
                    if (curEpochValidatorsIdMap[val] != 0) continue;
                    uint256 deposit = _validators[val].votes +
                        _validators[val].deposit;
                    if (randDeposit <= deposit) {
                        curEpochValidators.push(val);
                        curEpochValidatorsIdMap[val] = curEpochValidators
                            .length;
                        total -= deposit;
                        sysRewards.updateValidatorWhileElect(
                            val,
                            _validators[val].rate,
                            nextEpoch
                        );
                        break;
                    }
                    randDeposit -= deposit;
                }
            }
        }

        sysRewards.updateEpochWhileElect(
            totalTemp,
            curEpochValidators.length,
            effictiveLen,
            nextEpoch
        );
    }

    function rand(uint256 _length, uint256 _i) internal view returns (uint256) {
        uint256 random = uint256(
            keccak256(abi.encodePacked(blockhash(block.number - _i - 1), _i))
        );
        return random % _length;
    }

    function recentFourteenEpochAvgValCount() internal view returns (uint256) {
        uint256 curEpoch = currentEpoch();
        if (curEpoch == 0) {
            return effictiveValidators.length();
        }
        uint256 sumValidatorCount = 0;
        uint256 avg = 14;
        if (curEpoch < avg - 1) {
            avg = curEpoch;
        }
        for (uint256 i = 0; i < avg; i++) {
            (, , , uint256 effValCount) = sysRewards.epochs(curEpoch - i);
            sumValidatorCount += effValCount;
        }
        return sumValidatorCount / avg;
    }

    function nextEpochValidatorCount() internal view returns (uint256) {
        uint256 avgCount = recentFourteenEpochAvgValCount();
        if (avgCount < MIN_LEVEL_VALIDATOR_COUNT) {
            return MAX_VALIDATOR_COUNT_LV1;
        }
        if (avgCount < MEDIUM_LEVEL_VALIDATOR_COUNT) {
            return MAX_VALIDATOR_COUNT_LV2;
        }
        if (avgCount < MAX_LEVEL_VALIDATOR_COUNT) {
            return MAX_VALIDATOR_COUNT_LV3;
        }
        // avgCount >= MAX_LEVEL_VALIDATOR_COUNT
        return MAX_VALIDATOR_COUNT_LV4;
    }

}

interface IValidators {
    function isEffictiveValidator(address addr) external view returns (bool);

    function getEffictiveValidators() external view returns (address[] memory);

    function getInvalidValidators() external view returns (address[] memory);

    function effictiveValsLength() external view returns (uint256);

    function invalidValsLength() external view returns (uint256);

    function validators(address _val)
        external
        view
        returns (Validators.Validator calldata);

    function kickoutValidator(address _val) external;

    function tryElect() external;

    function addValidatorFromProposal(
        address _addr,
        uint256 _deposit,
        uint8 _rate,
        string memory _name,
        string memory _details
    ) external payable;

    function voteValidator(
        address _voter,
        address _val,
        uint256 _votes
    ) external payable;

    function cancelVoteValidator(
        address _voter,
        address _val,
        uint256 _votes,
        bool _clear
    ) external payable;

}

contract Proposals is Base, Initializable, ReentrancyGuard {
    using EnumerableSet for EnumerableSet.Bytes32Set;

    enum ProposalType {
        init
    }

    enum ProposalStatus {
        pending,
        pass,
        cancel
    }

    struct ProposalInfo {
        /// @notice id
        bytes4 id;
        address proposer;
        ProposalType pType;
        uint256 deposit;
        uint8 rate;
        /// @notice name
        string name;
        string details;
        uint256 initBlock;
        address guarantee;
        uint256 updateBlock;
        ProposalStatus status;
    }

    mapping(bytes4 => ProposalInfo) public proposalInfos;
    mapping(address => bytes4[]) public proposals;

    EnumerableSet.Bytes32Set proposalsBytes;

    /// @notice Validators contract
    IValidators public validators;

    event LogInitProposal(
        bytes32 indexed id,
        address indexed proposer,
        uint256 block,
        uint256 deposit,
        uint256 rate
    );
    event LogGuarantee(
        bytes32 indexed id,
        address indexed guarantee,
        uint256 block
    );
    event LogCancelProposal(
        bytes32 indexed id,
        address indexed proposer,
        uint256 block
    );
    event LogUpdateProposal(
        bytes32 indexed id,
        address indexed proposer,
        uint256 block,
        uint256 deposit,
        uint256 rate
    );

    modifier onlyEffictiveValidator() {
        require(
            validators.isEffictiveValidator(msg.sender) ||
                validators.effictiveValsLength() == 0,
            "Proposals: msg sender must be validator"
        );
        _;
    }

    modifier onlyEffictiveProposal(bytes4 id) {
        require(
            block.number <=
                proposalInfos[id].initBlock +
                    PROPOSAL_DURATION_EPOCHS *
                    EPOCH_BLOCKS,
            "Proposals: Proposal has expired"
        );
        _;
    }

    modifier checkValidatorLength() {
        require(
            validators.effictiveValsLength() < MAX_VALIDATORS_COUNT,
            "Proposals: length of the validator must be less than MAX_VALIDATORS_COUNT"
        );
        _;
    }

    function initialize(address _validator) public onlySystem initializer {
        validators = IValidators(_validator);
    }

    /**
     * @dev initProposal
     */
    function initProposal(
        ProposalType pType,
        uint8 rate,
        string memory name,
        string memory details
    ) external payable nonReentrant checkValidatorLength {
        require(
            !validators.isEffictiveValidator(msg.sender),
            "Proposals: The msg.sender can not be validator"
        );
        require(
            !Address.isContract(msg.sender),
            "Proposals: The msg.sender can not be contract address"
        );
        require(
            bytes(details).length <= MAX_VALIDATOR_DETAIL_LENGTH,
            "Proposals: Details is too long"
        );
        require(
            bytes(name).length <= MAX_VALIDATOR_NAME_LENGTH,
            "Proposals: name is too long"
        );
        require(
            msg.value >= MIN_DEPOSIT,
            "Proposals: Deposit must greater than MIN_DEPOSIT"
        );
        require(
            rate >= MIN_RATE && rate <= MAX_RATE,
            "Proposals: Rate must greater than MIN_RATE and less than MAX_RATE"
        );
        bytes4[] memory lastIds = proposals[msg.sender];
        if (lastIds.length > 0) {
            bytes4 lastId = lastIds[lastIds.length - 1];
            require(
                proposalInfos[lastId].status != ProposalStatus.pending,
                "Proposals: The msg.sender's latest proposal is still in pending"
            );
        }
        bytes4 id = bytes4(
            keccak256(
                abi.encodePacked(
                    msg.sender,
                    msg.value,
                    rate,
                    name,
                    details,
                    block.number
                )
            )
        );
        require(
            proposalInfos[id].initBlock == 0,
            "Proposals: Proposal already exists"
        );
        ProposalInfo memory proposal;
        proposal.deposit = msg.value;
        proposal.id = id;
        proposal.details = details;
        proposal.name = name;
        proposal.initBlock = block.number;
        proposal.proposer = msg.sender;
        proposal.pType = pType;
        proposal.status = ProposalStatus.pending;
        proposal.rate = rate;
        proposalInfos[id] = proposal;
        proposals[address(msg.sender)].push(id);
        proposalsBytes.add(id);
        emit LogInitProposal(id, msg.sender, block.number, msg.value, rate);
    }

    /**
     * @dev guarantee
     */
    function guarantee(bytes4 id)
        external
        nonReentrant
        onlyEffictiveValidator
        onlyEffictiveProposal(id)
    {
        require(
            proposalInfos[id].initBlock != 0,
            "Proposals: proposal not exist"
        );
        require(
            proposalInfos[id].status == ProposalStatus.pending,
            "Proposals: The status of proposal must be pending"
        );
        proposalInfos[id].updateBlock = block.number;
        proposalInfos[id].guarantee = msg.sender;
        validators.addValidatorFromProposal{value: proposalInfos[id].deposit}(
            proposalInfos[id].proposer,
            proposalInfos[id].deposit,
            proposalInfos[id].rate,
            proposalInfos[id].name,
            proposalInfos[id].details
        );
        proposalInfos[id].status = ProposalStatus.pass;
        emit LogGuarantee(id, msg.sender, block.number);
    }

    /**
     * @dev updateProposal
     */
    function updateProposal(
        bytes4 id,
        uint8 rate,
        uint256 deposit,
        string memory name,
        string memory details
    ) external payable nonReentrant onlyEffictiveProposal(id) {
        require(
            proposalInfos[id].initBlock != 0,
            "Proposals: proposal not exist"
        );
        require(
            proposalInfos[id].proposer == msg.sender,
            "Proposals: not proposer"
        );
        require(
            proposalInfos[id].status == ProposalStatus.pending,
            "Proposals: The status of proposal must be pending"
        );
        require(
            bytes(name).length <= MAX_VALIDATOR_NAME_LENGTH,
            "Proposals: name is too long"
        );
        require(
            bytes(details).length <= MAX_VALIDATOR_DETAIL_LENGTH,
            "Proposals: details is too long"
        );
        require(
            deposit >= MIN_DEPOSIT,
            "Proposals: deposit must greater than MIN_DEPOSIT"
        );
        require(
            rate >= MIN_RATE && rate <= MAX_RATE,
            "Proposals: rate must greater than MIN_RATE and less than MAX_RATE"
        );
        uint256 lastDeposit = proposalInfos[id].deposit;
        if (lastDeposit > deposit) {
            address payable receiver = payable(address(msg.sender));
            receiver.transfer(lastDeposit - deposit);
        } else if (lastDeposit < deposit) {
            require(
                deposit - lastDeposit == msg.value,
                "Proposals: msg value not true"
            );
        } else {
            if (msg.value != 0) {
                address payable receiver = payable(address(msg.sender));
                receiver.transfer(msg.value);
            }
        }
        proposalInfos[id].deposit = deposit;
        proposalInfos[id].rate = rate;
        proposalInfos[id].updateBlock = block.number;
        proposalInfos[id].name = name;
        proposalInfos[id].details = details;
        emit LogUpdateProposal(id, msg.sender, block.number, deposit, rate);
    }

    /**
     * @dev cancelProposal
     */
    function cancelProposal(bytes4 id) external nonReentrant {
        require(
            proposalInfos[id].initBlock != 0,
            "Proposals: proposal not exist"
        );
        require(
            proposalInfos[id].proposer == msg.sender,
            "Proposals: not proposer"
        );
        require(
            proposalInfos[id].status == ProposalStatus.pending,
            "Proposals: The status of proposal must be pending"
        );
        proposalInfos[id].updateBlock = block.number;
        address payable receiver = payable(address(msg.sender));
        receiver.transfer(proposalInfos[id].deposit);
        proposalInfos[id].status = ProposalStatus.cancel;
        emit LogCancelProposal(id, msg.sender, block.number);
    }

    function allProposals(uint256 page, uint256 size)
        public
        view
        returns (ProposalInfo[] memory)
    {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        uint256 start = (page - 1) * size;
        if (proposalsBytes.length() < start) {
            size = 0;
        } else {
            uint256 length = proposalsBytes.length() - start;
            if (length < size) {
                size = length;
            }
        }

        ProposalInfo[] memory proposalDir = new ProposalInfo[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = proposalInfos[
                bytes4(proposalsBytes.at(i + start))
            ];
        }
        return proposalDir;
    }

    function allProposalSets(uint256 page, uint256 size)
        public
        view
        returns (bytes4[] memory)
    {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        uint256 start = (page - 1) * size;
        if (proposalsBytes.length() < start) {
            size = 0;
        } else {
            uint256 length = proposalsBytes.length() - start;
            if (length < size) {
                size = length;
            }
        }
        bytes4[] memory proposalDir = new bytes4[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = bytes4(proposalsBytes.at(i + start));
        }
        return proposalDir;
    }

    function addressProposals(
        address val,
        uint256 page,
        uint256 size
    ) public view returns (ProposalInfo[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        bytes4[] memory addressProposalIds = proposals[val];
        uint256 start = (page - 1) * size;
        if (addressProposalIds.length < start) {
            size = 0;
        } else {
            uint256 length = addressProposalIds.length - start;
            if (length < size) {
                size = length;
            }
        }

        ProposalInfo[] memory proposalDir = new ProposalInfo[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = proposalInfos[addressProposalIds[i + start]];
        }
        return proposalDir;
    }

    function addressProposalSets(
        address val,
        uint256 page,
        uint256 size
    ) public view returns (bytes4[] memory) {
        require(page > 0 && size > 0, "Proposals: Requests param error");
        bytes4[] memory addressProposalIds = proposals[val];
        uint256 start = (page - 1) * size;
        if (addressProposalIds.length < start) {
            size = 0;
        } else {
            uint256 length = addressProposalIds.length - start;
            if (length < size) {
                size = length;
            }
        }

        bytes4[] memory proposalDir = new bytes4[](size);
        for (uint256 i = 0; i < size; i++) {
            proposalDir[i] = addressProposalIds[i + start];
        }
        return proposalDir;
    }

    function proposalCount() public view returns (uint256) {
        return proposalsBytes.length();
    }

    function addressProposalCount(address val) public view returns (uint256) {
        return proposals[val].length;
    }

}