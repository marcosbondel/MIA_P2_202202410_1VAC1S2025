
export const appReducer = ( state, action) => {
    switch (action.type) {
        case 'disks[set]':
            return {
                ...state,
                disks: action.payload.disks
            }
        case 'error[set]':
            console.log("action.payload: ", action.payload)
            console.log("action.payload.error: ", action.payload['error'])
            return {
                ...state,
                showError: action.payload.error != "",
                error: action.payload.error
            }
        default:
            break;
    }

}