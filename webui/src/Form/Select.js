import React from "react";
import styled from "styled-components";
import MSelect from "@mui/material/Select";
import InputLabel from "@mui/material/InputLabel";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";

// Styled Component
const FullWidthFormControl = styled(FormControl)`
  && {
    min-width: 100%;
  }
`;

const render = ({ label, id, options, ...props }) => (
  <FullWidthFormControl>
    <InputLabel htmlFor={id}>{label}</InputLabel>
    <MSelect inputProps={{ name: id, id }} {...props}>
      {options.map((o) => (
        <MenuItem key={o.name} value={o.val}>
          {o.name}
        </MenuItem>
      ))}
    </MSelect>
  </FullWidthFormControl>
);

export default render;
