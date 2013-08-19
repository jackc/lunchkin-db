require 'capybara/rspec'
require 'sequel'
require 'pry'
require 'yaml'

config = YAML.load File.read('config.test.yml')

DB = Sequel.postgres database: config['database']['database']

Capybara.default_driver = :selenium
Capybara.app_host = "http://#{config['address']}:#{config['port']}"
