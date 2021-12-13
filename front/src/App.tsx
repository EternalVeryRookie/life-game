import { useState, useCallback, useEffect } from "react";

import { LifeGameWorld, useLifeGame } from "./api_client";
import { genRandomOrderArray } from "./rand";
import { Box, Slider, Button } from "@mui/material";

const lifeStateBaseStyle = {width: "5px", height: "5px"}
const Alive = <Box sx={{backgroundColor: "blue", ...lifeStateBaseStyle}}/>
const Death = <Box sx={{backgroundColor: "red" , ...lifeStateBaseStyle}}/>

function WorldGrid({world}: {world: LifeGameWorld}) {
  return (
    <>{
      <Box>{
        world.map((row, i) => {
          return <Box sx={{display: "flex"}} key={i}>{row.map((cell, j) => cell === "alive" ? Alive : Death)}</Box>
        })
      }</Box>
    }</>
  );
}

function genRandomLifeGameWorld(width: number, height: number, aliveNum: number): LifeGameWorld {
  if (aliveNum > width*height) {
    aliveNum = width * height;
  }

  const deathNum = width*height - aliveNum;
  const arr = genRandomOrderArray([{value: "alive" as "alive", num: aliveNum}, {value: "death" as "death", num: deathNum}]);
  return Array(height).fill(null).map((_, i) => arr.slice(i*width, width+i*width));
}

const width = 64;
const height = 64;

function App() {
  const [world, setWorld] = useState(() => genRandomLifeGameWorld(width, height, width*height/2));
  const methods = useLifeGame(setWorld);
  const [isSimulationing, setIsSimulationing] = useState(false);

  const onChange = useCallback((_, newValue) => {
    const newWorld = genRandomLifeGameWorld(width, height, newValue);
    setWorld(newWorld);
  }, [setWorld]);
  
  const onClickStart = useCallback(() => {
    methods?.start(world);
    setIsSimulationing(true);
  }, [world, methods]);

  const onClickStop = useCallback(() => methods?.stop(), [methods]);
  const onClickEnd = useCallback(() =>{
    methods?.end();
    setIsSimulationing(false);
  }, [methods]);

  useEffect(() => methods?.next());

  return (
    <Box>
      <WorldGrid world={world}/>
      <Slider disabled={isSimulationing} onChange={onChange} defaultValue={width*height/2} min={0} max={width*height}/>
      <Button onClick={onClickStart}>Start</Button>
      <Button onClick={onClickStop}>Stop</Button>
      <Button onClick={onClickEnd}>End</Button>
    </Box>
  )
}

export default App;
