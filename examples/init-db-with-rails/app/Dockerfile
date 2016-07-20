
FROM    rails:5

WORKDIR /code
COPY    setup.sh /code/
RUN     ./setup.sh

WORKDIR /code/blog
CMD     ["bin/rails", "server", "-b", "0.0.0.0"]
