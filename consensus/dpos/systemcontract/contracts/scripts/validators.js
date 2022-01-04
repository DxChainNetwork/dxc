
const { web3, validators } = require("./contracts.js")

const getCurEpochValidators = async () => {
    const curEpochValidators = await validators.methods.getCurEpochValidators().call()
    console.log(curEpochValidators);
}

getCurEpochValidators()
