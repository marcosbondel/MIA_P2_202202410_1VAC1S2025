import { Breadcrumbs, Container, Toolbar, Typography, Link, Grid, Button, IconButton } from "@mui/material"
import { useContext, useEffect, useState } from "react";
import { FileCard } from "../components/FileCard";
import { AppContext } from "../../context/AppContext";
import { useNavigate } from "react-router-dom";
import { Link as RouterLink } from "react-router-dom";
import { Document } from "../components/Document";

export const PartitionPage = () => {
    const { 
        getFileSystem, 
        current_directory, 
        directory_parts, 
        current_fs_location,
        dispatchCurrentFSLocation,
        current_file_content
    } = useContext(AppContext)


    const navigate = useNavigate()


    const handleOnClick = (e, part) => {
        e.preventDefault();
        // Handle click event for breadcrumb links
        // This could be used to navigate to the parent directory or a specific path
        console.log(`Clicked part: ${part}`);
        let previousPath = current_fs_location.split(part).slice(0, -1).join('/');
        console.log(`Navigating to: ${previousPath}`);
        dispatchCurrentFSLocation(`${previousPath}${part}`);
    }

    const handleGoBack = () => {
        // Handle the go back action
        // This could be used to navigate to the previous directory or a specific path
        const parts = current_fs_location.split('/');
        parts.pop(); // Remove the last part to go back one level
        let previousPath = parts.join('/');
        
        if(previousPath == '')
            previousPath = '/';
        console.log(`Going back to: ${previousPath}`);

        dispatchCurrentFSLocation(previousPath);
    }

    useEffect(() => {
        getFileSystem()
    }, [current_fs_location]);


    return (
        <Container>
            <Toolbar/>
            <h1>Partition Page</h1>
            { current_fs_location && current_fs_location != '/' ? (

                <Breadcrumbs aria-label="breadcrumb" sx={{ marginBottom: '1rem', backgroundColor: '#C0C9EE', padding: '1rem', color: '#fff', borderRadius: '4px' }}>
                    { directory_parts && ( directory_parts.map(part => (
                        <Link underline="hover" color="inherit" href="/" onClick={(e) => handleOnClick(e, part)}>
                            {part}
                        </Link>
                    )))}
                </Breadcrumbs>
            ) : (
                <Breadcrumbs aria-label="breadcrumb" sx={{ marginBottom: '1rem', backgroundColor: '#C0C9EE', padding: '1rem', color: '#fff', borderRadius: '4px' }}>
                    <Link underline="hover" color="inherit" >
                        /
                    </Link>
                </Breadcrumbs>
            )}
            
            <Toolbar>
                {/* <Link href="/mia" onClick={() => navigate('/mia')}>Go back</Link> */}
                {/* <Link component={RouterLink} color='inherit' to='/auth/login'>Go back</Link> */}
                { current_fs_location !== '/' ? (
                    <Button variant="outlined" onClick={handleGoBack}>Go back</Button>
                ):(
                    <Button variant="outlined" disabled>Go back</Button>
                )}
                {/* <IconButton>
                    Go back
                </IconButton> */}
            </Toolbar>
            <Grid
                container
                spacing={6}
                direction="row"
                alignItems="center"
                sx={{ padding: '1rem' }}
            >
                { (current_directory && current_directory.length > 0 ? (
                    current_directory.map((dir, index) => (
                        <FileCard key={`${dir.name}-${index}`} dir={dir} path={current_fs_location}/>
                    ))
                ) : (
                    // <Typography>Nothing to show</Typography>
                    <Document content={current_file_content} />
                )) }
            </Grid>
        </Container>
    )
}
