require 'capybara/rspec'
require 'sequel'
require 'pry'

DB = Sequel.postgres database: 'munchdb_development'

Capybara.default_driver = :selenium
Capybara.app_host = 'http://localhost:4000'
