import { Breadcrumbs, Container, Toolbar, Typography, Link, Grid } from "@mui/material"
import { useEffect, useState } from "react";
import { DocumentCard } from "../components/DocumentCard";
// import {FolderIcon, ArticleIcon} from '@mui/icons-material';

export const PartitionPage = () => {
    const [current_fs, set_current_fs] = useState({})
    const [current_directory, set_current_directory] = useState([])
    const [directory_parts, setDirectory_parts] = useState([])


    useEffect(() => {
        fetch("http://localhost:3000/api/fs?path=/home/user/docs")
            .then((res) => res.json())
            .then((data) => {
                console.log(data)
                set_current_fs(data);
                set_current_directory(data?.children);
                setDirectory_parts(data.path.split("/").slice(1));
                if (data.type === "directory") {
                    console.log("Es carpeta:", data.children);
                } else {
                    console.log("Contenido archivo:", data.content);
                }
        });
    }, []);


    return (
        <Container>
            <Toolbar/>
            <h1>Partition Page</h1>
            <Breadcrumbs aria-label="breadcrumb">
                { directory_parts && directory_parts.length > 0 ? ( directory_parts.map(part => (
                    <Link underline="hover" color="inherit" href="/">
                        {part}
                    </Link>
                ))): (
                    <Typography color="text.primary">No directories to show</Typography>
                ) }
            </Breadcrumbs>
            <Grid
                container
                spacing={2}
                direction="row"
                justifyContent="space-between"
                alignItems="center"
            >
                { (current_directory && current_directory.length > 0 ? (
                    current_directory.map((dir, index) => (
                        <DocumentCard key={`${dir.name}-${index}`} dir={dir} />
                    ))
                ) : (
                    <Typography>Nothing to show</Typography>
                )) }
            </Grid>
            <pre>
                current_directory: {JSON.stringify(current_directory, null, 2)}<br/>
            </pre>
        </Container>
    )
}
