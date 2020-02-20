FROM hashicorp/terraform:0.12.10

WORKDIR /work
COPY . /work

ENV GOOGLE_APPLICATION_CREDENTIALS "/work/terraform-key.json"
RUN terraform init
ENTRYPOINT ["terraform"]
