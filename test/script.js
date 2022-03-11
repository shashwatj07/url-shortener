// K6 HTTP load generator

import http from "k6/http";
// import { check } from 'k6';

export default function () {
  const url =
    "http://url-shortener3-dev.ap-south-1.elasticbeanstalk.com/f";
  // const payload = JSON.stringify({
  //   urldata: "https://www.google.com",
  // });

  // const params = {
  //   headers: {
  //     "Content-Type": "application/json",
  //   },
  // };

  const res = http.get(url);

  // console.log(res.body);
}
