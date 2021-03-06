require 'capybara/rspec'
require 'sequel'
require 'pry'
require 'yaml'

Dir["#{File.dirname(__FILE__)}/support/**/*.rb"].each {|f| require f}

config = YAML.load File.read('config.test.yml')

DB = Sequel.postgres database: config['database']['database']

Capybara.default_driver = :selenium
Capybara.app_host = "http://#{config['address']}:#{config['port']}"

RSpec.configure do |config|
  config.before(:each) do
    clean_database
    visit '/'
  end
end