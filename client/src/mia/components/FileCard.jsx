import { Box, Card, CardActions, CardContent, Grid, IconButton, Typography } from "@mui/material"
import { useLocation, useNavigate, useParams } from "react-router-dom";
import ArticleOutlined from '@mui/icons-material/ArticleOutlined';
import queryString from "query-string"
import FolderOutlined from '@mui/icons-material/FolderOutlined';
import { useContext } from "react";
import { AppContext } from "../../context/AppContext";

export const FileCard = ({dir}) => {

    const { dispatchCurrentFSLocation, current_fs_location } = useContext(AppContext);

    const navigate = useNavigate();
    const { disk, partition } = useParams()
    const location = useLocation()
    const { path = '' } = queryString.parse(location.path)

    const handleDocumentClick = () => {
        // navigate(`/mia/disks/${disk}/partitions/${partition}?path=${path}/${dir.name}`);
        if (current_fs_location === '/') {
            dispatchCurrentFSLocation(`${current_fs_location}${dir.name}`)
        } else {
            dispatchCurrentFSLocation(`${current_fs_location}/${dir.name}`)
        }
    }

    return (
        <Grid item xs={12} sm={6} md={4} lg={3}>
        {/* <CardContent> */}
            {/* <Box> */}
                <IconButton onClick={handleDocumentClick}>
                    { dir.type === "directory" ? (
                        <FolderOutlined sx={{ fontSize: '6rem', color: 'primary.main' }}/>
                    ) : (
                        <ArticleOutlined sx={{ fontSize: '6rem', color: 'primary.main' }}/>
                    )}
                </IconButton>
                <Typography variant="h5" align="center">{ dir.name }</Typography>
            {/* </Box> */}
        {/* </CardContent> */}
        </Grid>
    )
}
