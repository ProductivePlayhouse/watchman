import React from 'react';
import { StyleSheetManager } from 'styled-components';

const StylesProvider = ({ children }) =>
{
    const nonce = window.cspNonce; // This is set by the script in index.html

    return (
        <StyleSheetManager nonce={nonce}>
            {children}
        </StyleSheetManager>
    );
};

export default StylesProvider;