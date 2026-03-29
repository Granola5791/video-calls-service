import { useState } from "react"
import { General, Auth, Errors } from "../constants/hebrew-constants"
import { CenteredColumn, CenteredFilledScreen } from "../styled-components/StyledBoxes"
import { LongButtonFilled } from "../styled-components/StyledButtons"
import { LoadingText, StyledTitle } from "../styled-components/StyledText"
import { StyledTextField } from "../styled-components/StyledTextFields"
import { ErrorText } from "../styled-components/StyledErrors"
import { UsersServer, HttpStatusCodes } from "../constants/backend-constants"
import { useNavigation } from "../utils/navigation"
import { Typography } from "@mui/material"
import { Link } from "react-router-dom"
import { LocalStorage, RouterPaths } from "../constants/general-contants"
import { StyledLogo, StyledLogoContainer } from "../styled-components/StyledLogo"

const LoginPage = () => {

    const [response, setResponse] = useState<string>('');
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [isError, setIsError] = useState<boolean>(true);

    const { goToHome: GoToHome } = useNavigation();

    const handleLogin = async (username: string, password: string) => {
        setIsError(false);
        setResponse(Auth.wait);
        try {
            const res = await fetch(UsersServer.httpAddress + UsersServer.api.logIn, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password }),
                credentials: 'include'
            })
            switch (res.status) {
                case HttpStatusCodes.OK:
                    const data = await res.json();
                    localStorage.setItem(LocalStorage.id, data.id);
                    localStorage.setItem(LocalStorage.role, data.role);
                    localStorage.setItem(LocalStorage.username, data.username);
                    GoToHome();
                    break;
                case HttpStatusCodes.Unauthorized:
                    setIsError(true);
                    setResponse(Errors.invalidAuthInput);
                    break;
                default:
                    setIsError(true);
                    setResponse(Errors.genericError);
            }
        } catch (error) {
            setIsError(true);
            setResponse(Errors.genericError);
        }
    }


    return (

        <CenteredFilledScreen>
            <CenteredColumn >
                <StyledLogoContainer>
                    <StyledLogo src="/assets/logo.jpg" alt="Logo" />
                    <StyledTitle>{General.appName}</StyledTitle>
                </StyledLogoContainer>

                <CenteredColumn>
                    <Typography>
                        {Auth.noAccountYet.text}
                        <Link to={RouterPaths.signup}>{Auth.noAccountYet.linkText}</Link>
                    </Typography>
                    <StyledTextField
                        label={Auth.username}
                        variant="filled"
                        onChange={(e) => setUsername(e.target.value)}
                    />
                    <StyledTextField
                        label={Auth.password}
                        type="password"
                        variant="filled"
                        onChange={(e) => setPassword(e.target.value)}
                    />

                    <LongButtonFilled onClick={() => handleLogin(username, password)}>{Auth.loginButton}</LongButtonFilled>

                    {isError ?
                        <ErrorText>{response}</ErrorText>
                        :
                        <LoadingText>{response}</LoadingText>
                    }
                </ CenteredColumn>
            </CenteredColumn>
        </CenteredFilledScreen>
    )
}

export default LoginPage