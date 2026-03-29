import styled from "@emotion/styled";
import { Colors } from "../constants/general-contants";
export const StyledLogo = styled.img({
    width: 'auto',
    height: '80px',
});


export const StyledLogoContainer = styled.div({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: 'auto',
    marginBottom: '0px',
    flexDirection: 'row',
    gap: '10px',
    backgroundColor: 'white',
    padding: '0px',
    borderRadius: '10px',
});

export const CornerLogoContainer = styled(StyledLogoContainer)({
    position: 'absolute',
    top: '12px',
    right: '20px',
    cursor: 'pointer',
});

export const CornerLogo = styled(StyledLogo)({
    height: '60px',
});

export const LogoTitle = styled.h1({
    fontSize: '3rem',
    color: Colors.primary,
    margin: '0',
});