import { Typography } from "@mui/material"

export const Document = ({content}) => {
    return (
        <Typography variant="p" sx={{width: '100%', backgroundColor: "#FBFBFB", padding: '2rem'}} >
            {content }
        </Typography>
    )
}
