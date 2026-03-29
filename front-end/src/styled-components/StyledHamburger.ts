import { styled } from "@mui/material";

export const HamburgerDrawerContainer = styled('div')(({ theme }) => ({
    position: "fixed",
    width: "20vw",
    height: "100vh",
    border: "1px solid black",
    backgroundColor: theme.palette.background.paper,
}));

export const HamburgerMenuTitle = styled('h1')(({ theme }) => ({
    textAlign: 'center',
    color: theme.palette.primary.main,
}));