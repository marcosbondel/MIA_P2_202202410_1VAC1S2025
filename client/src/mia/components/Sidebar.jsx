import { TurnedInNot } from "@mui/icons-material"
import { Box, Divider, Drawer, Grid, List, ListItem, ListItemButton, ListItemIcon, ListItemText, Toolbar, Typography } from "@mui/material"

export const Sidebar = ({drawerWidth}) => {
    return (
        <Drawer
            variant="permanent"
            open
            sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
        >
            <List>
                {
                    ['Enero', 'Febrero', 'Marzo', 'Abril'].map(text => (
                        <ListItem key={text} disablePadding>
                            <ListItemButton>
                                <ListItemIcon>
                                    <TurnedInNot/>
                                </ListItemIcon>
                                <Grid container>
                                    <ListItemText primary={text}></ListItemText>
                                    {/* <ListItemText secondary={'Lorem ipsum, dolor sit amet.'}></ListItemText> */}
                                </Grid>
                            </ListItemButton>
                        </ListItem>
                    ))
                }
            </List>
        </Drawer>
    )
}
