import { Box, Card, CardActions, CardContent, Grid, IconButton, Typography } from "@mui/material"
import { useNavigate, useParams } from "react-router-dom";
import ArticleIcon from '@mui/icons-material/ArticleOutlined';

export const DocumentCard = ({dir}) => {

    const navigate = useNavigate();
    const { disk } = useParams()

    return (
        <Grid item xs={12} sm={6} md={4} lg={3}>
        {/* <CardContent> */}
            <Box>
                <IconButton>
                    <ArticleIcon sx={{ fontSize: '6rem', color: 'primary.main' }}/>
                </IconButton>
                <Typography variant="h5" align="center">{ dir.name }</Typography>
            </Box>
        {/* </CardContent> */}
        </Grid>
    )
}
