import styled from "@emotion/styled";
import { Colors } from "../constants/general-contants";


export const CenteredColumn = styled.div({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: 'auto',
    flexDirection: 'column',
    gap: '10px',
    backgroundColor: 'white',
    padding: '20px',
    borderRadius: '10px',
})

export const CenteredRow = styled.div({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: 'auto',
    flexDirection: 'row',
    gap: '10px',
    backgroundColor: 'white',
    padding: '20px',
    borderRadius: '10px',
})

export const CenteredScreen = styled.div({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    flexDirection: 'column',
    height: '100vh',
    width: '100vw',
    margin: 'auto',
    gap: '20px',
})

export const CenteredFilledScreen = styled(CenteredScreen)({
    backgroundColor: Colors.primary,
})

export const StyledMeetingGrid = styled.div({
    display: "grid",
    position: "fixed",
    gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))",
    gap: "8px",
    width: "100%",
    height: "87vh",
    top: "0px",
    right: "0px",
    left: "0px",
    bottom: "0px",
    padding: "8px",
    boxSizing: "border-box",
    backgroundColor: Colors.meetingBackground,
    overflow: "hidden",
});

export const HamburgerMenuContainer = styled.div({
    position: "fixed",
    width: "20vw",
    height: "100vh",
    border: "1px solid black",
    backgroundColor: Colors.primary,
});