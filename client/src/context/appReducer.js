
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
        case 'logged[set]':
            return {
                ...state,
                logged: action.payload.logged
            }
        case 'partitions[set]':
            return {
                ...state,
                partitions: action.payload.partitions
            }
        default:
            break;
    }

}