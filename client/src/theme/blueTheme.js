import { createTheme } from "@mui/material";
import { red } from "@mui/material/colors";

export const blueTheme = createTheme({
    palette: {
        primary: {
            main: '#000957'
        },
        secondary: {
            main: '#344CB7'
        },
        error: {
            main: red.A400
        }
    }
})