import React from "react";
import styled from "styled-components";
import MSelect from "@mui/material/Select";
import InputLabel from "@mui/material/InputLabel";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";

const render = ({ label, id, options, ...props }) => (
  <FormControl
    css={`
      && {
        min-width: 100%;
      }
    `}
  >
    <InputLabel htmlFor={id}>{label}</InputLabel>
    <MSelect inputProps={{ name: id, id }} {...props}>
      {options.map(o => (
        <MenuItem key={o.name} value={o.val}>
          {o.name}
        </MenuItem>
      ))}
    </MSelect>
  </FormControl>
);

export default render;
