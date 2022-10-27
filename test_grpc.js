import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['filmepb'], 'filme.proto');

export const options = {
  discardResponseBodies: true,
  scenarios: {
    contacts: {
      executor: 'constant-vus',
      vus: 1,
      duration: '20s',
    },
  },
};

export default () => {
  client.connect('localhost:50052', {
    plaintext:true,
  })

  const data = { };
  const response = client.invoke('filme.MoviesService/ListarFilmesGrpc', data);

  // check(response, {
  //   'status is OK': (r) => r && r.status === grpc.StatusOK,
  // });

  // console.log(JSON.stringify(response.message));

  client.close();
  sleep(1)
};