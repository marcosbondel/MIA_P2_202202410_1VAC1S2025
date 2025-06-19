import { ThemeProvider } from "@emotion/react"
import { CssBaseline } from "@mui/material"
import { blueTheme } from "./blueTheme"

export const AppTheme = ({children}) => {
    return (
        <ThemeProvider theme={blueTheme}>
            <CssBaseline/>
            { children }
        </ThemeProvider>
    )
}
