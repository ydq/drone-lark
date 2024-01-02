FROM istio/distroless:1.20-2023-12-20T19-00-54
MAINTAINER darren <i@darren.work>
COPY ./lark /lark
ENTRYPOINT ["/lark"]