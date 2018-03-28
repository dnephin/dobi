#!/bin/bash
#
# Automate the steps of the rails getting started guide
#
set -eu

rails new blog
cd blog
bin/rails generate controller Welcome index

cat > config/routes.rb <<EOF
Rails.application.routes.draw do

  resources :articles

  root 'welcome#index'
end
EOF

bin/rails generate controller Articles
cat > app/controllers/articles_controller.rb <<EOF
class ArticlesController < ApplicationController
  def new
  end

  def create
    @article = Article.new(article_params)

    @article.save
    redirect_to @article
  end

  private
    def article_params
      params.require(:article).permit(:title, :text)
    end
end
EOF

cat > app/views/articles/new.html.erb <<EOF
<%= form_for :article, url: articles_path do |f| %>
  <p>
    <%= f.label :title %><br>
    <%= f.text_field :title %>
  </p>

  <p>
    <%= f.label :text %><br>
    <%= f.text_area :text %>
  </p>

  <p>
    <%= f.submit %>
  </p>
<% end %>
EOF

bin/rails generate model Article title:string text:text

# TODO: set these to different values
cat > config/database.yml <<EOF
development:
  adapter: postgresql
  database: postgres
  host: postgres
  user: postgres
EOF

echo "gem 'pg', '~> 0.18.0'" >> ./Gemfile
bundle install
