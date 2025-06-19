
export const appReducer = ( state = initial_state, action) => {
    switch (action.type) {
        case 'options[add]':
            
            return {
                ...state,
                options: action.payload.options
            }
        case 'result[add]':
            return {
                ...state,
                error: {},
                result: action.payload
            }
            
        case 'error[add]':
            return {
                ...state,
                result: {},
                error: action.payload
            }
        default:
            break;
    }

}