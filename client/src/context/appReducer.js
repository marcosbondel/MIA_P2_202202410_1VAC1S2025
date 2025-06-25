
export const appReducer = ( state, action) => {
    switch (action.type) {
        case 'disks[set]':
            return {
                ...state,
                disks: action.payload.disks
            }
        case 'error[set]':
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
        case 'current_fs_location[set]':
            return {
                ...state,
                current_fs_location: action.payload.current_fs_location
            }
        case 'current_fs[set]':
            return {
                ...state,
                current_fs: action.payload.current_fs
            }
        case 'current_directory[set]':
            return {
                ...state,
                current_directory: action.payload.current_directory
            }
        case 'directory_parts[set]':
            return {
                ...state,
                directory_parts: action.payload.directory_parts
            }
        case 'current_file_content[set]':
            return {
                ...state,
                current_file_content: action.payload.current_file_content
            }
        default:
            break;
    }

}