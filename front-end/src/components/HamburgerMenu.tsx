import { useState } from "react";
import { Divider, Drawer, List, ListItem, ListItemButton } from '@mui/material';
import { HamburgerMenuButton } from "../styled-components/StyledButtons";
import { HamburgerMenuIcon } from "../styled-components/StyledIcons";
import type { MenuOption } from '../types/menuOptions'
import { HamburgerDrawerContainer, HamburgerMenuTitle } from "../styled-components/StyledHamburger";
import { StyledP } from "../styled-components/StyledText";

interface HamburgerMenuProps {
    onClose?: () => void,
    onOpen?: () => void,
    topButtons?: MenuOption[],
    bottomButtons?: MenuOption[],
    title?: string
    closeOnClick?: boolean;
}

const HamburgerMenu = ({ onClose, onOpen, topButtons = [], bottomButtons = [], title, closeOnClick }: HamburgerMenuProps) => {
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

    const handleButtonClick = (button: MenuOption) => {
        if (closeOnClick) {
            handleClose();
        }
        button.onClick();
    };

    return (
        <>
            <HamburgerMenuButton onClick={handleOpen}>
                <HamburgerMenuIcon />
            </HamburgerMenuButton>
            <Drawer anchor='left' open={open} onClose={handleClose}>
                <HamburgerDrawerContainer>
                    <HamburgerMenuTitle>{title}</HamburgerMenuTitle>
                    <List>
                        {topButtons.map((button, index) => (
                            <ListItem key={index} disablePadding>
                                <ListItemButton onClick={() => handleButtonClick(button)}>
                                    <StyledP>
                                        {button.text}
                                    </StyledP>
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                    <Divider />
                    <List>
                        {bottomButtons.map((button, index) => (
                            <ListItem key={index} disablePadding>
                                <ListItemButton onClick={() => handleButtonClick(button)}>
                                    <StyledP>
                                        {button.text}
                                    </StyledP>
                                </ListItemButton>
                            </ListItem>
                        ))}
                    </List>
                </HamburgerDrawerContainer>
            </Drawer>
        </>
    )
}

export default HamburgerMenu