import { General, Auth } from "../constants/hebrew-constants"
import { CenteredColumn } from "../styled-components/StyledBoxes"
import { LongButtonFilled } from "../styled-components/StyledButtons"
import { StyledTitle } from "../styled-components/StyledText"
import { StyledTextField } from "../styled-components/StyledTextFields"

const LoginPage = () => {
    return (
        <>
            <StyledTitle>{General.appName}</StyledTitle>

            <CenteredColumn>
                <StyledTextField
                    label={Auth.username}
                    variant="filled"
                />
                <StyledTextField
                    label={Auth.password}
                    type="password"
                    variant="filled"
                />

                <LongButtonFilled>{Auth.loginButton}</LongButtonFilled>
            </ CenteredColumn>

        </>
    )
}

export default LoginPage