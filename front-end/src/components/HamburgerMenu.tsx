import { useState } from "react";
import { Divider, Drawer, List, ListItem, ListItemButton } from '@mui/material';
import { HamburgerMenuButton } from "../styled-components/StyledButtons";
import { HamburgerMenuContainer } from "../styled-components/StyledBoxes";
import { HamburgerMenuIcon } from "../styled-components/StyledIcons";

interface HamburgerMenuProps {
    onClose?: () => void,
    onOpen?: () => void,
    topButtons?: { text: string, onClick: () => void }[],
    bottomButtons?: { text: string, onClick: () => void }[],
    title?: string
}

const HamburgerMenu = ({ onClose, onOpen, topButtons, bottomButtons, title}: HamburgerMenuProps) => {

    const [open, setOpen] = useState(false);

    const handleOpen = () => {
        setOpen(true);
        if (onOpen) {
            onOpen();
        }
    };

    const handleClose = () => {
        setOpen(false);
        if (onClose) {
            onClose();
        }
    };

    return (
        <div style={{width: 'fit-content', border: '1px solid black'}}>
            <HamburgerMenuButton onClick={handleOpen}>
                <HamburgerMenuIcon />
            </HamburgerMenuButton>
            <Drawer anchor='right' open={open} onClose={handleClose}>
                <HamburgerMenuContainer>
                    {title && <h4 style={{textAlign: 'center'}}>{title}</h4>}
                    <List>
                        {topButtons && topButtons.map((button, index) => (
                            <ListItem key={index} disablePadding>
                                <ListItemButton onClick={button.onClick}>
                                    {button.text}
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                    <Divider />
                    <List>
                        {bottomButtons && bottomButtons.map((button, index) => (
                            <ListItem key={index} disablePadding>
                                <ListItemButton onClick={button.onClick}>
                                    {button.text}
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                </HamburgerMenuContainer>
            </Drawer>
        </div>
    )
}

export default HamburgerMenu