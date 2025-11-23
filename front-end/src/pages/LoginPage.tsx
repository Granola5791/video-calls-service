import { useState } from "react"
import { General, Auth, Errors, SuccessMessages } from "../constants/hebrew-constants"
import { CenteredColumn, CenteredFilledScreen } from "../styled-components/StyledBoxes"
import { LongButtonFilled } from "../styled-components/StyledButtons"
import { StyledTitle } from "../styled-components/StyledText"
import { StyledTextField } from "../styled-components/StyledTextFields"
import { ErrorText } from "../styled-components/StyledErrors"
import { ApiEndpoints, BackendAddress, HttpStatusCodes } from "../constants/backend-constants"
import { useNavigation } from "../utils/navigation"

const LoginPage = () => {

    const [response, setResponse] = useState<string>('');
    const [username, setUsername] = useState<string>('');
    const [password, setPassword] = useState<string>('');

    const { goToHome: GoToHome } = useNavigation();

    const handleLogin = async (username: string, password: string) => {
        setResponse(Auth.wait);
        const res = await fetch(BackendAddress + ApiEndpoints.logIn, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
            credentials: 'include'
        })
        switch (res.status) {
            case HttpStatusCodes.OK:
                GoToHome();
                break;
            case HttpStatusCodes.Unauthorized:
                setResponse(Errors.invalidAuthInput);
                break;
            default:
                setResponse(Errors.genericError);
        }
    }


    return (

        <CenteredFilledScreen>
            <CenteredColumn >
                <StyledTitle>{General.appName}</StyledTitle>

                <CenteredColumn>
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

                    <ErrorText>{response}</ErrorText>
                </ CenteredColumn>
            </CenteredColumn>
        </CenteredFilledScreen>
    )
}

export default LoginPage