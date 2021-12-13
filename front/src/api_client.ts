import React from "react";

const api_url = "ws://lifegame-env.eba-mymxdspc.us-east-2.elasticbeanstalk.com/";

export type LifeState = "alive" | "death";
export type LifeGameWorld = LifeState[][];
type methods = {
  start: (world: LifeGameWorld) => void,
  next: () => void,
  stop: () => void, 
  end: () => void,
}

export function useLifeGame(onRecieveWorldState: (world: LifeGameWorld) => void) {
  const [methods, setMethods] = React.useState<methods|null>(null);

  React.useEffect(() => {
    const conn = new WebSocket(api_url + "lifegame");
 
    conn.onopen = () => {
      const onmessage = (message: MessageEvent<string>) => {
        const world = JSON.parse(message.data) as LifeGameWorld;
        onRecieveWorldState(world);
      }
      conn.onmessage = onmessage;

      const start = (world: LifeGameWorld) => conn.send("start"+JSON.stringify(world));
      const next = () => conn.send("next");
      const stop = () => conn.send("stop");
      const end = () => conn.send("end");
      setMethods({
        start,next,stop,end
      });
    }
  

    return () => conn.close();
  }, [onRecieveWorldState]);

  return methods;
}