import styled from "@emotion/styled";
import { Colors } from "../constants/general-contants";

export const HamburgerDrawerContainer = styled.div({
    position: "fixed",
    width: "20vw",
    height: "100vh",
    border: "1px solid black",
    backgroundColor: Colors.drawerBackground,
});

export const HamburgerMenuTitle = styled.h4({
    textAlign: 'center',
    color: Colors.primary,
});